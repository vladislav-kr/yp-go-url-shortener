package app

import (
	"net/http"
	"time"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/handlers"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/router"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/server"
	mapkeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper"
)

type Option struct {
	Host            string
	RedirectHost    string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

func NewServer(opt Option) *server.HTTPServer {

	// Обработчики с доступом к хранилищу
	h := handlers.NewHandlers(
		mapkeeper.New(),
		opt.RedirectHost,
	)

	srv := &http.Server{
		Addr:         opt.Host,
		Handler:      router.NewRouter(h),
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
		IdleTimeout:  opt.IdleTimeout,
	}

	return &server.HTTPServer{
		Server:          srv,
	}
}
