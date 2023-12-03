package save

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/keeper"
)

func TestNewSaveHandler(t *testing.T) {
	const redirectHost = "http://localhost:8080"

	stor := keeper.New()

	tests := []struct {
		name           string
		sorage         keeper.Keeperer
		url            string
		expectedStatus int
	}{
		{
			name:           "url ok",
			sorage:         stor,
			url:            "https://ya.ru/",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "url empty",
			sorage:         stor,
			url:            "",
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(tt.url))
			require.NoError(t, err)

			NewSaveHandler(tt.sorage, redirectHost).ServeHTTP(rr, req)

			result := rr.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.expectedStatus, result.StatusCode)
		})
	}
}
