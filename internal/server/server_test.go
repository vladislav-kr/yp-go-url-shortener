package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"golang.org/x/sync/errgroup"
)

func TestHTTPServer(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	srv := NewHTTPServer(
		zaptest.NewLogger(t),
		ts.Config,
	)

	errgroup, errgroupContext := errgroup.WithContext(ctx)

	errgroup.Go(func() error {
		return srv.Run()
	})

	errgroup.Go(func() error {
		<-errgroupContext.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return srv.Stop(ctx)
	})

	assert.NoError(t, errgroup.Wait())

}
