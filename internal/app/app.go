// app сборка и запуск сервера
package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/handlers"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/middleware"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/middleware/auth"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/router"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/lib/fileutils"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/server"
	urlHandler "github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler/deleter"
	dbkeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/db-keeper"
	mapkeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper"
)

// URLShortener хранить все параметры сервера.
type URLShortener struct {
	log             *zap.Logger
	server          *server.HTTPServer
	shutdownTimeout time.Duration
	db              *sql.DB
	memStorage      *mapkeeper.Keeper
}

// Option конфигурация сервера.
type Option struct {
	Host            string
	RedirectHost    string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	StorageFilePath string
	StorageDBDNS    string
}

// NewURLShortener новая инстанция сервера.
func NewURLShortener(ctx context.Context, log *zap.Logger, opt Option) (*URLShortener, error) {
	var (
		db         *sql.DB
		memStorage *mapkeeper.Keeper
		storage    urlHandler.Keeperer
	)

	switch {
	case len(opt.StorageDBDNS) > 0:
		var err error
		db, err = connectionDB(opt.StorageDBDNS)
		if err != nil {
			return nil, err
		}

		dbPool, err := connectionPool(ctx, opt.StorageDBDNS)
		if err != nil {
			return nil, err
		}

		storage = dbkeeper.NewDBKeeper(log.With(
			zap.String(
				"component",
				"dbkeeper",
			),
		),
			db,
			dbPool,
		)
	case len(opt.StorageFilePath) > 0:
		storageFilePath, err := validateStorageFilePath(opt.StorageFilePath)
		if err != nil {
			return nil, err
		}
		memStorage = mapkeeper.New(storageFilePath)
		storage = memStorage
	default:
		return nil, fmt.Errorf("failed to create storage")
	}

	h := handlers.NewHandlers(
		log.With(
			zap.String(
				"component",
				"handlers",
			),
		),
		urlHandler.NewURLHandler(
			storage,
			db,
			deleter.NewDeleter(ctx, 10, func(urls []models.DeleteURL) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				storage.DeleteURLS(ctx, urls)
			})),
		opt.RedirectHost,
	)

	m := middleware.New(
		log.With(
			zap.String(
				"component",
				"middleware",
			),
		),
		auth.New("secret-key"),
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
		shutdownTimeout: opt.ShutdownTimeout,
		db:              db,
		memStorage:      memStorage,
	}, nil
}

// Run запускает сервер.
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
		if us.db != nil {

			tx, err := us.db.BeginTx(ctx, nil)
			if err != nil {
				us.log.Info(
					"failed to create transaction",
					zap.Error(err),
				)
				return err
			}

			defer tx.Rollback()
			_, err = tx.ExecContext(ctx, `
			CREATE TABLE
				IF NOT EXISTS shortened_url (
					short_url VARCHAR(10) PRIMARY KEY,
					original_url VARCHAR(4000) NOT NULL
				);
			`)
			if err != nil {
				us.log.Info(
					"failed to create table shortened_url",
					zap.Error(err),
				)
				return err
			}
			_, err = tx.ExecContext(ctx,
				`CREATE UNIQUE INDEX IF NOT EXISTS orig_url_idx ON shortened_url (original_url)`)
			if err != nil {
				us.log.Info(
					"failed to create index",
					zap.String("field", "original_url"),
					zap.Error(err),
				)
				return err
			}

			_, err = tx.ExecContext(ctx,
				`ALTER TABLE shortened_url ADD COLUMN IF NOT EXISTS user_id UUID;`)
			if err != nil {
				us.log.Info(
					"failed to create new column user_id",
					zap.Error(err),
				)
				return err
			}

			_, err = tx.ExecContext(ctx,
				`ALTER TABLE shortened_url ADD COLUMN IF NOT EXISTS is_deleted boolean DEFAULT FALSE;`)
			if err != nil {
				us.log.Info(
					"failed to create new column is_deleted",
					zap.Error(err),
				)
				return err
			}

			if err := tx.Commit(); err != nil {
				us.log.Info(
					"failed to apply changes to the database",
					zap.Error(err),
				)
				return err
			}
		}

		if us.memStorage != nil {
			if err := us.memStorage.LoadFromFile(); err != nil {
				us.log.Info(
					"failed to load data from file",
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
			if us.memStorage != nil {
				if err := us.memStorage.SaveToFile(); err != nil {
					us.log.Error(
						"failed to save data to file",
						zap.Error(err),
					)
				}
			}

		}()

		defer func() {
			if us.db != nil {
				if err := us.db.Close(); err != nil {
					us.log.Error(
						"error closing connection to database",
						zap.Error(err),
					)
				}
			}
		}()

		return us.server.Stop(ctx)
	})

	return errGr.Wait()

}

func connectionDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, fmt.Errorf("failed connection to database: %w", err)
	}
	return db, nil
}

func connectionPool(ctx context.Context, dns string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, dns)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	return db, nil
}

func validateStorageFilePath(path string) (string, error) {
	storageFilePath, err := fileutils.CreateFullPathFromRelative(path)
	if err != nil {
		return "", err
	}

	dir, _ := filepath.Split(storageFilePath)

	if _, err := os.Stat(dir); err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			if err := os.Mkdir(dir, os.ModeDir); err != nil {
				return "", err
			}
		default:
			return "", err
		}

	}
	return storageFilePath, nil
}
