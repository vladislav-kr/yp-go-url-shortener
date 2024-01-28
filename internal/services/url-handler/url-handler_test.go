package urlhandler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler/mocks"
)

func TestReadURL(t *testing.T) {

	cases := []struct {
		name        string
		id          string
		expectedURL string
		expectedErr error
		isError     bool
		isCallMock  bool
	}{
		{
			name:        "successful receiving url",
			id:          "idurltest",
			expectedURL: "https://practicum.yandex.ru/",
			isCallMock:  true,
		},
		{
			name:        "failed to receive url",
			id:          "idurltest",
			expectedErr: errors.New("not found url"),
			isError:     true,
			isCallMock:  true,
		},
		{
			name:        "alias is empty",
			expectedErr: errors.New("alias is empty"),
			isError:     true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			storage := mocks.NewKeeperer(t)

			if tc.isCallMock {
				storage.On("GetURL", mock.AnythingOfType("*context.timerCtx"), tc.id).
					Return(tc.expectedURL, tc.expectedErr)
			}

			h := NewURLHandler(
				storage,
				mocks.NewDBPinger(t),
			)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			url, err := h.ReadURL(ctx, tc.id)

			if tc.isError {
				assert.Empty(t, url)
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedURL, url)

		})

	}

}

func TestSaveURL(t *testing.T) {

	cases := []struct {
		name          string
		longURL       string
		expectedAlias string
		expectedErr   error
		isError       bool
		isCallMock    bool
	}{
		{
			name:          "successful receiving url",
			longURL:       "https://practicum.yandex.ru/",
			expectedAlias: "test-alias",
			isCallMock:    true,
		},
		{
			name:        "failed to receive url",
			longURL:     "https://practicum.yandex.ru/",
			expectedErr: errors.New("not found url"),
			isError:     true,
			isCallMock:  true,
		},
		{
			name:        "invalid url",
			longURL:     "practicum.yandex.ru/",
			expectedErr: errors.New("invalid url"),
			isError:     true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			storage := mocks.NewKeeperer(t)

			if tc.isCallMock {
				storage.On("PostURL", mock.AnythingOfType("*context.timerCtx"), tc.longURL, "").
					Return(tc.expectedAlias, tc.expectedErr)
			}

			h := NewURLHandler(storage, mocks.NewDBPinger(t))
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			alias, err := h.SaveURL(ctx, tc.longURL, "")

			if tc.isError {
				assert.Empty(t, alias)
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedAlias, alias)

		})

	}

}

func TestPing(t *testing.T) {

	cases := []struct {
		name    string
		isError bool
	}{
		{
			name: "successful database ping",
		},
		{
			name:    "unsuccessful database ping",
			isError: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			storage := mocks.NewKeeperer(t)

			var err error
			if tc.isError {
				err = errors.New("fail ping db")
			}

			pingDB := mocks.NewDBPinger(t)

			h := NewURLHandler(storage, pingDB)
			pingDB.
				On("PingContext", mock.AnythingOfType("*context.timerCtx")).
				Return(err)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			errPing := h.Ping(ctx)

			if tc.isError {
				assert.Error(t, errPing)
				return
			}

			assert.NoError(t, errPing)

		})

	}

}

func TestSaveURLS(t *testing.T) {

	cases := []struct {
		name         string
		urls         []models.BatchRequest
		expectedURLS []models.BatchResponse
		err          error
		isError      bool
	}{
		{
			name: "successful receiving url",
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
					ShortURL:      "https://localhost:8080/dkh2ksukde",
				},
				{
					CorrelationID: "2",
					ShortURL:      "https://localhost:8080/fh43jfhfdq",
				},
			},
		},
		{
			name:         "failed to receive url",
			urls:         []models.BatchRequest{},
			expectedURLS: nil,
			err:          errors.New("no data to save"),
			isError:      true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			storage := mocks.NewKeeperer(t)

			storage.On("SaveURLS", mock.AnythingOfType("*context.timerCtx"), tc.urls, "").
				Return(tc.expectedURLS, tc.err)

			h := NewURLHandler(storage, mocks.NewDBPinger(t))
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			respURLS, err := h.SaveURLS(ctx, tc.urls, "")

			if tc.isError {
				assert.Nil(t, respURLS)
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedURLS, respURLS)

		})

	}

}
