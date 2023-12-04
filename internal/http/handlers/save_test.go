package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mapkeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper"
)

func TestSaveHandler(t *testing.T) {

	h := NewHandlers(mapkeeper.New(), "http://localhost:8080")

	tests := []struct {
		name           string
		handler        func(w http.ResponseWriter, r *http.Request)
		url            string
		expectedStatus int
	}{
		{
			name:           "url ok",
			handler:        h.SaveHandler,
			url:            "https://ya.ru/",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "url empty",
			handler:        h.SaveHandler,
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

			tt.handler(rr, req)

			result := rr.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.expectedStatus, result.StatusCode)
		})
	}
}
