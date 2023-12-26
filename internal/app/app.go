package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/handlers"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/middleware"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/router"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/lib/fileutils"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/server"
	urlHandler "github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler"
	mapkeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type URLShortener struct {
	log             *zap.Logger
	server          *server.HTTPServer
	storage         *mapkeeper.Keeper
	storageFilePath string
	shutdownTimeout time.Duration
}

type Option struct {
	Host            string
	RedirectHost    string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	StorageFilePath string
}

func NewURLShortener(log *zap.Logger, opt Option) (*URLShortener, error) {

	storageFilePath, err := fileutils.CreateFullPathFromRelative(opt.StorageFilePath)
	if err != nil {
		return nil, err
	}

	dir, _ := filepath.Split(storageFilePath)

	if err := os.Mkdir(dir, os.ModeDir); err != nil &&
		!errors.Is(err, syscall.ERROR_ALREADY_EXISTS) {
		return nil, err
	}

	memStorage := mapkeeper.New()

	h := handlers.NewHandlers(
		log.With(
			zap.String(
				"component",
				"handlers",
			),
		),
		urlHandler.NewURLHandler(
			memStorage,
		),
		opt.RedirectHost,
	)

	m := middleware.New(
		log.With(
			zap.String(
				"component",
				"middleware",
			),
		),
	)

	srv := &http.Server{
		Addr:         opt.Host,
		Handler:      router.NewRouter(h, m),
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
		IdleTimeout:  opt.IdleTimeout,
	}

	return &URLShortener{
		log: log,
		server: server.NewHTTPServer(
			log.With(
				zap.String("component", "HTTPServer"),
				zap.String("addr", srv.Addr),
			),
			srv),
		storage:         memStorage,
		storageFilePath: storageFilePath,
		shutdownTimeout: opt.ShutdownTimeout,
	}, nil
}

func (us *URLShortener) Run(ctx context.Context) error {
	// Контекст прослушивающий сигналы прерывания OS
	sigCtx, sigCancel := signal.NotifyContext(ctx,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer sigCancel()

	// Группа для запуска и остановки сервера по сигналу
	errGr, errGrCtx := errgroup.WithContext(sigCtx)

	errGr.Go(func() error {
		if len(us.storageFilePath) > 0 {
			if err := us.storage.LoadFromFile(us.storageFilePath); err != nil {
				us.log.Info(
					"failed to load data from file",
					zap.String("path", us.storageFilePath),
					zap.Error(err),
				)
			}
		}

		return us.server.Run()
	})

	errGr.Go(func() error {
		<-errGrCtx.Done()

		ctx, cancel := context.WithTimeout(
			context.Background(),
			us.shutdownTimeout,
		)
		defer cancel()

		defer func() {
			if err := us.storage.SaveToFile(us.storageFilePath); err != nil {
				us.log.Error(
					"failed to save data to file",
					zap.String("path", us.storageFilePath),
					zap.Error(err),
				)
			}
		}()

		return us.server.Stop(ctx)
	})

	return errGr.Wait()

}
