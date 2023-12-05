package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mapkeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper"
)

func TestRedirectHandler(t *testing.T) {
	stor := mapkeeper.New()

	id, err := stor.PostURL("https://ya.ru/")
	require.NoError(t, err)
	h := NewHandlers(stor, "http://localhost:8080")

	tests := []struct {
		name             string
		handler          func(w http.ResponseWriter, r *http.Request)
		id               string
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "positive test",
			handler:          h.RedirectHandler,
			id:               id,
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: "https://ya.ru/",
		},
		{
			name:           "negative test",
			handler:        h.RedirectHandler,
			id:             "no-id",
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()

			r.Use(middleware.URLFormat)

			r.Get("/{id}", tt.handler)

			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", ts.URL, tt.id), nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedLocation, resp.Header.Get("Location"))
		})
	}
}
