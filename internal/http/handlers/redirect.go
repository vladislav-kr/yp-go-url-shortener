package handlers

import (
	"net/http"
	"strings"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/keeper"
)

func NewRedirectHandler(sorage keeper.Keeperer) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/")
		url, err := sorage.GetURL(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h := w.Header()
		h.Set("Location", url)
		h.Set("Content-Type", "text/html; charset=utf-8")

		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
