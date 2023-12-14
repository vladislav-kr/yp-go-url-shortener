package server

import (
	"context"
	"errors"
	"net/http"
)

type HTTPServer struct {
	Server *http.Server
}

func (hs *HTTPServer) Run() error {
	if err := hs.Server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (hs *HTTPServer) Stop(ctx context.Context) error {
	return hs.Server.Shutdown(ctx)
}
