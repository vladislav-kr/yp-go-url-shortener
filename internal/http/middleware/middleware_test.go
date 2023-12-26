package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestLogger(t *testing.T) {
	log := zaptest.NewLogger(
		t,
		zaptest.WrapOptions(zap.Hooks(func(e zapcore.Entry) error {
			assert.Equal(t, e.Message, "request completed")
			return nil
		})),
	)

	r := chi.NewRouter()
	r.Use(New(log).Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/", ts.URL), nil)
	require.NoError(t, err)

	client := &http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

}
