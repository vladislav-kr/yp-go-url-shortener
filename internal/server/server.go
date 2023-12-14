package server

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type HTTPServer struct {
	Server          *http.Server
	ShutdownTimeout time.Duration
}

func (hs *HTTPServer) Run() error {
	if err := hs.Server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (hs *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), hs.ShutdownTimeout)
	defer cancel()

	return hs.Server.Shutdown(ctx)
}
