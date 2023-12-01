package handlers

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
	stor := keeper.New()

	tests := []struct {
		name           string
		sorage         keeper.Keeperer
		url            string
		expectedStatus int
	}{
		{
			name:           "positive test",
			sorage:         stor,
			url:            "https://ya.ru/",
			expectedStatus: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(tt.url))
			require.NoError(t, err)

			NewSaveHandler(tt.sorage).ServeHTTP(rr, req)

			result := rr.Result()

			assert.Equal(t, tt.expectedStatus, result.StatusCode)
		})
	}
}
