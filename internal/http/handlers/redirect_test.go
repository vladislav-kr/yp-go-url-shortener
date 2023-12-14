package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/keeper"
)

func TestNewRedirectHandler(t *testing.T) {

	stor := keeper.New()
	id, err := stor.PostURL("https://ya.ru/")
	require.NoError(t, err)

	tests := []struct {
		name             string
		sorage           keeper.Keeperer
		id               string
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "positive test",
			sorage:           stor,
			id:               id,
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: "https://ya.ru/",
		},
		{
			name:           "negative test",
			sorage:         stor,
			id:             "no-id",
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.id), nil)
			require.NoError(t, err)

			NewRedirectHandler(tt.sorage).ServeHTTP(rr, req)

			result := rr.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.expectedStatus, result.StatusCode)
			assert.Equal(t, tt.expectedLocation, result.Header.Get("Location"))
		})
	}
}
