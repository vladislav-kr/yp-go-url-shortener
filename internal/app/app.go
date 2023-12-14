package app

import (
	"net/http"
	"time"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/handlers"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/server"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/keeper"
)

func NewServer(host string) *server.HTTPServer {

	mux := http.NewServeMux()

	stor := keeper.New()
	mux.HandleFunc("/", handlers.NewURLHandler(stor))

	srv := &http.Server{
		Addr:    host,
		Handler: mux,
	}

	return &server.HTTPServer{
		Server:          srv,
		ShutdownTimeout: 10 * time.Second,
	}
}
