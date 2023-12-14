package save

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/keeper"
)

func NewSaveHandler(
	sorage keeper.Keeperer,
	redirectHost string,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := sorage.PostURL(string(data))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)

		render.PlainText(w, r, fmt.Sprintf("%s/%s", redirectHost, id))
	}
}
