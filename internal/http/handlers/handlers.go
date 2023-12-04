package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Keeperer interface {
	PostURL(url string) (string, error)
	GetURL(id string) (string, error)
}

type Handlers struct {
	storage      Keeperer
	redirectHost string
}

func NewHandlers(
	storage Keeperer,
	redirectHost string,
) *Handlers {

	return &Handlers{
		storage:      storage,
		redirectHost: redirectHost,
	}
}

func (h *Handlers) SaveHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.storage.PostURL(string(data))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)

	render.PlainText(w, r, fmt.Sprintf("%s/%s", h.redirectHost, id))
}

func (h *Handlers) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := h.storage.GetURL(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.WriteHeader(http.StatusTemporaryRedirect)
}
