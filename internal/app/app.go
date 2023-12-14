package app

import (
	"net/http"
	"time"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/router"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/server"
)

type Option struct {
	Host string
	RedirectHost string
}

func NewServer(opt Option) *server.HTTPServer {
	mux, _ := router.NewRouter(opt.RedirectHost)

	srv := &http.Server{
		Addr:    opt.Host,
		Handler: mux,
	}

	return &server.HTTPServer{
		Server:          srv,
		ShutdownTimeout: 10 * time.Second,
	}
}
