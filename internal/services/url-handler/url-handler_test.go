package urlhandler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
				storage.On("GetURL", tc.id).
					Return(tc.expectedURL, tc.expectedErr)
			}

			h := NewURLHandler(storage)

			url, err := h.ReadURL(tc.id)

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
				storage.On("PostURL", tc.longURL).
					Return(tc.expectedAlias, tc.expectedErr)
			}

			h := NewURLHandler(storage)

			alias, err := h.SaveURL(tc.longURL)

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
