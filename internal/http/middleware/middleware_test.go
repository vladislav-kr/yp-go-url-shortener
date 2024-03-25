package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
	r.Use(New(log, nil).Logger)
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

func testCompress(t *testing.T, data []byte) []byte {
	var b bytes.Buffer

	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		t.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	return b.Bytes()
}
func TestCompress(t *testing.T) {
	log := zaptest.NewLogger(t)

	cases := []struct {
		name                  string
		headers               map[string]string
		ContentTypesSupported []string
		responseContentType   string
		requestBody           string
		successBody           string
	}{
		{
			name: "compressed json data",
			headers: map[string]string{
				"Content-Encoding": "gzip",
				"Accept-Encoding":  "gzip",
			},
			ContentTypesSupported: []string{
				"application/json",
			},
			responseContentType: "application/json",
			requestBody:         `{"url": "https://practicum.yandex.ru"}`,
			successBody:         `{"result": "http://localhost:8080/sjdhsjdjhsd"}`,
		},
		{
			name: "no compressed data",
			headers: map[string]string{
				"Accept-Encoding": "gzip",
			},
			ContentTypesSupported: []string{
				"text/plain",
			},
			responseContentType: "text/plain; charset=utf-8",
			requestBody:         "Request data",
			successBody:         "Response data",
		},
		{
			name: "compressed data",
			headers: map[string]string{
				"Content-Encoding": "gzip",
				"Accept-Encoding":  "gzip",
			},
			ContentTypesSupported: []string{
				"text/plain",
			},
			responseContentType: "text/plain; charset=utf-8",
			requestBody:         "Request data",
			successBody:         "Response data",
		},
		{
			name: "compressed data with an unsupported content type",
			headers: map[string]string{
				"Content-Encoding": "gzip",
				"Accept-Encoding":  "gzip",
			},
			ContentTypesSupported: []string{},
			responseContentType:   "text/plain; charset=utf-8",
			requestBody:           "Request data",
			successBody:           "Response data",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := chi.NewRouter()
			r.Use(New(log, nil).NewCompressHandler(tc.ContentTypesSupported))
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				b, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				if strings.Compare(string(b), tc.requestBody) == 0 {
					w.Header().Set("Content-Type", tc.responseContentType)
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tc.successBody))
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
			})

			ts := httptest.NewServer(r)
			defer ts.Close()

			body := bytes.NewBuffer(nil)
			if tc.headers["Content-Encoding"] == "gzip" {
				body.Write(testCompress(t, []byte(tc.requestBody)))
			} else {
				body.Write([]byte(tc.requestBody))
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/", ts.URL), body)
			require.NoError(t, err)

			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			var b []byte

			if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
				zr, err := gzip.NewReader(resp.Body)
				require.NoError(t, err)

				b, err = io.ReadAll(zr)
				require.NoError(t, err)
			} else {
				b, err = io.ReadAll(resp.Body)
				require.NoError(t, err)
			}
			require.Equal(t, tc.successBody, string(b))

		})

	}

}
