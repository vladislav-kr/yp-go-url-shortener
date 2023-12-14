package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/keeper"
)

func NewSaveHandler(sorage keeper.Keeperer) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
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
		w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", id)))
		
	}
}