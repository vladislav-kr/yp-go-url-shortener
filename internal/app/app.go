package app

import (
	"net/http"
	"time"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/handlers"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/middleware"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/router"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/server"
	urlHandler "github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler"
	mapkeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper"
	"go.uber.org/zap"
)

type Option struct {
	Host         string
	RedirectHost string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func NewServer(log *zap.Logger, opt Option) *server.HTTPServer {

	h := handlers.NewHandlers(
		log.With(
			zap.String(
				"component",
				"handlers",
			),
		),
		urlHandler.NewURLHandler(
			mapkeeper.New(),
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

	return server.NewHTTPServer(
		log.With(
			zap.String(
				"component",
				"HTTPServer",
			),
			zap.String(
				"addr",
				srv.Addr,
			),
		),
		srv)
}
