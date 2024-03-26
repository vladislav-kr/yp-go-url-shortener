// server отвечает за запуск и остановку http сервера
package server

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// HTTPServer хранит информацию о http сервере.
type HTTPServer struct {
	log    *zap.Logger
	server *http.Server
}

// NewHTTPServer новый http-сервер.
//
//	NewHTTPServer(zap.L(),&http.Server{Addr:":8080"})
func NewHTTPServer(
	log *zap.Logger,
	srv *http.Server,
) *HTTPServer {
	return &HTTPServer{
		log:    log,
		server: srv,
	}
}

// Run запускает сервер.
func (hs *HTTPServer) Run() error {
	hs.log.Info("running")
	err := hs.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		hs.log.Warn("stopped")
		return nil
	} else {
		hs.log.Error(
			"failed to stopped",
			zap.Error(err),
		)
		return err
	}
}

// Stop graceful shutdown сервера.
func (hs *HTTPServer) Stop(ctx context.Context) error {
	hs.log.Info("stopping...")
	return hs.server.Shutdown(ctx)
}
