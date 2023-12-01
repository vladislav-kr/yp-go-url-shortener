package server

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type HttpServer struct {
	Server          *http.Server
	ShutdownTimeout time.Duration
}

func (hs *HttpServer) Run() error {
	if err := hs.Server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (hs *HttpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), hs.ShutdownTimeout)
	defer cancel()

	return hs.Server.Shutdown(ctx)
}
