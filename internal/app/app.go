package app

import (
	"net/http"
	"time"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/router"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/server"
)

func NewServer(host string) *server.HTTPServer {
	mux, _ := router.NewRouter()

	srv := &http.Server{
		Addr:    host,
		Handler: mux,
	}

	return &server.HTTPServer{
		Server:          srv,
		ShutdownTimeout: 10 * time.Second,
	}
}
