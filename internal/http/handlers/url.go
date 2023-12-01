package handlers

import (
	"net/http"
	"strings"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/keeper"
)

func NewURLHandler(sorage keeper.Keeperer) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := strings.TrimPrefix(r.URL.Path, "/")

		if len(id) != 0 {
			NewRedirectHandler(sorage).ServeHTTP(w, r)
			return	
		}
		NewSaveHandler(sorage).ServeHTTP(w, r)
	}
}
