package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestHTTPServer(t *testing.T) {

	t.Run("start and stop server", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		srv := HTTPServer{
			Server:          ts.Config,
			ShutdownTimeout: time.Second,
		}

		errgroup, errgroupContext := errgroup.WithContext(ctx)

		errgroup.Go(func() error {
			return srv.Run()
		})

		errgroup.Go(func() error {
			<-errgroupContext.Done()
			return srv.Stop()
		})

		assert.NoError(t, errgroup.Wait())
	})

}
