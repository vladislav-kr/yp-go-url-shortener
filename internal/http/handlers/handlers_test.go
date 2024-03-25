package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/handlers/mocks"
	urlhandler "github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler/deleter"
	mapkeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper"
)

func TestSaveHandler(t *testing.T) {

	cases := []struct {
		name           string
		url            string
		alias          string
		err            error
		expectedStatus int
	}{
		{
			name:           "url ok",
			url:            "https://ya.ru/",
			alias:          "alias1",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "url empty",
			url:            "",
			err:            errors.New("url empty"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid url",
			url:            "ya.ru",
			err:            errors.New("invalid url"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			urlHndl := mocks.NewURLHandler(t)
			urlHndl.On("SaveURL", mock.AnythingOfType("*context.timerCtx"), tc.url, "").
				Return(tc.alias, tc.err)

			h := NewHandlers(
				zaptest.NewLogger(t),
				urlHndl,
				"http://localhost:8080",
			)

			rr := httptest.NewRecorder()
			req, err := http.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(tc.url),
			)
			require.NoError(t, err)

			h.SaveHandler(rr, req)

			result := rr.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.expectedStatus, result.StatusCode)

		})
	}

}

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name             string
		alias            string
		err              error
		expectedStatus   int
		expectedLocation string
		isCallMock       bool
	}{
		{
			name:             "successful redirect",
			alias:            "alias1",
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: "https://ya.ru/",
			isCallMock:       true,
		},
		{
			name:           "unsuccessful redirect",
			alias:          "alias1",
			err:            errors.New("url not found"),
			expectedStatus: http.StatusBadRequest,
			isCallMock:     true,
		},
		{
			name:           "alias is empty: 404",
			err:            errors.New("alias is empty"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlHndl := mocks.NewURLHandler(t)
			if tc.isCallMock {
				urlHndl.On("ReadURL", mock.AnythingOfType("*context.timerCtx"), tc.alias).
					Return(tc.expectedLocation, tc.err)
			}

			h := NewHandlers(
				zaptest.NewLogger(t),
				urlHndl,
				"http://localhost:8080",
			)

			r := chi.NewRouter()

			r.Use(middleware.URLFormat)

			r.Get("/{id}", h.RedirectHandler)

			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", ts.URL, tc.alias), nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			assert.Equal(t, tc.expectedLocation, resp.Header.Get("Location"))

		})
	}
}

func TestSaveJSONHandler(t *testing.T) {

	cases := []struct {
		name           string
		url            models.URLRequest
		alias          string
		err            error
		expectedStatus int
	}{
		{
			name: "url ok",
			url: models.URLRequest{
				URL: "https://ya.ru/",
			},
			alias:          "alias1",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "url empty",
			err:            errors.New("url empty"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid url",
			url: models.URLRequest{
				URL: "ya.ru/",
			},
			err:            errors.New("invalid url"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			urlHndl := mocks.NewURLHandler(t)
			urlHndl.On("SaveURL", mock.AnythingOfType("*context.timerCtx"), tc.url.URL, "").
				Return(tc.alias, tc.err)

			h := NewHandlers(
				zaptest.NewLogger(t),
				urlHndl,
				"http://localhost:8080",
			)

			rr := httptest.NewRecorder()

			body, err := json.Marshal(tc.url)
			require.NoError(t, err)

			req, err := http.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(string(body)),
			)
			require.NoError(t, err)

			h.SaveJSONHandler(rr, req)

			result := rr.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.expectedStatus, result.StatusCode)

		})
	}

}

func TestPingHandler(t *testing.T) {

	cases := []struct {
		name           string
		expectedStatus int
		isError        bool
	}{
		{
			name:           "successful database ping",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "unsuccessful database ping",
			expectedStatus: http.StatusInternalServerError,
			isError:        true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			urlHndl := mocks.NewURLHandler(t)

			var err error
			if tc.isError {
				err = errors.New("fail ping db")
			}

			urlHndl.On("Ping", mock.AnythingOfType("*context.timerCtx")).
				Return(err)

			h := NewHandlers(
				zaptest.NewLogger(t),
				urlHndl,
				"http://localhost:8080",
			)

			rr := httptest.NewRecorder()

			req, err := http.NewRequest(
				http.MethodGet,
				"/",
				nil,
			)
			require.NoError(t, err)

			h.PingHandler(rr, req)

			result := rr.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.expectedStatus, result.StatusCode)

		})
	}

}

func TestBatchHandler(t *testing.T) {

	cases := []struct {
		name           string
		urls           []models.BatchRequest
		expectedURLS   []models.BatchResponse
		err            error
		expectedStatus int
		isError        bool
	}{
		{
			name: "data saved",
			urls: []models.BatchRequest{
				{
					CorrelationID: "1",
					OriginalURL:   "https://practicum.yandex.ru/",
				},
				{
					CorrelationID: "2",
					OriginalURL:   "https://ya.ru/",
				},
			},
			expectedURLS: []models.BatchResponse{
				{
					CorrelationID: "1",
					ShortURL:      "dkh2ksukde",
				},
				{
					CorrelationID: "2",
					ShortURL:      "fh43jfhfdq",
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "no data to save",
			urls:           []models.BatchRequest{},
			expectedURLS:   []models.BatchResponse{},
			err:            errors.New("no data to save"),
			expectedStatus: http.StatusBadRequest,
			isError:        true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			urlHndl := mocks.NewURLHandler(t)

			if !tc.isError {
				urlHndl.On("SaveURLS", mock.AnythingOfType("*context.timerCtx"), tc.urls, "").
					Return(tc.expectedURLS, tc.err)
			}

			h := NewHandlers(
				zaptest.NewLogger(t),
				urlHndl,
				"http://localhost:8080",
			)

			rr := httptest.NewRecorder()

			body, err := json.Marshal(tc.urls)
			require.NoError(t, err)

			req, err := http.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(string(body)),
			)
			require.NoError(t, err)

			h.BatchHandler(rr, req)

			result := rr.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.expectedStatus, result.StatusCode)

			if !tc.isError && result.ContentLength > 0 {
				respURLS := []models.BatchResponse{}
				err = json.NewDecoder(result.Body).Decode(&respURLS)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedURLS, respURLS)
			}

		})
	}

}

func ExampleHandlers_SaveHandler() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем in-memory хранилище
	storage := mapkeeper.New("")

	// Создаем обработчик сервисного слоя
	urlHandler := urlhandler.NewURLHandler(
		storage,
		nil,
		deleter.NewDeleter(ctx, 10, func(urls []models.DeleteURL) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			storage.DeleteURLS(ctx, urls)
		}),
	)

	// Создаем инстанцию http обработчиков
	h := NewHandlers(
		zap.L(),
		urlHandler,
		"http://localhost:8080",
	)

	router := chi.NewRouter()
	// Регистрируем обработчик на роутере
	router.Post("/", h.SaveHandler)

	ctxReq, cancelReq := context.WithTimeout(ctx, time.Second*2)
	defer cancelReq()

	// Создаем запрос на сокращение URL
	req, _ := http.NewRequestWithContext(
		ctxReq,
		http.MethodPost,
		"/",
		strings.NewReader("https://go.dev/"),
	)
	ww := httptest.NewRecorder()

	router.ServeHTTP(ww, req)
	ww.Result().Body.Close()
	fmt.Println(ww.Result().StatusCode)
	// Output:
	// 201

}
