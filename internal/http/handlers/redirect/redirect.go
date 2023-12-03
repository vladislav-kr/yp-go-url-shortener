package redirect

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/keeper"
)

func NewRedirectHandler(
	sorage keeper.Keeperer,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		if len(id) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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
