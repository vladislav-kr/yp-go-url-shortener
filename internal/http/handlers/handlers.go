package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	"go.uber.org/zap"
)

//go:generate mockery --name URLHandler
type URLHandler interface {
	ReadURL(alias string) (string, error)
	SaveURL(url string) (string, error)
}

type Handlers struct {
	log          *zap.Logger
	urlHandler   URLHandler
	redirectHost string
}

func NewHandlers(
	log *zap.Logger,
	urlHandler URLHandler,
	redirectHost string,
) *Handlers {
	return &Handlers{
		log:          log,
		urlHandler:   urlHandler,
		redirectHost: redirectHost,
	}
}

func (h *Handlers) SaveHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)

	if err != nil {
		h.log.Error(
			"failed reading body request",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.urlHandler.SaveURL(string(data))
	if err != nil {
		h.log.Error(
			"failed to save url",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, fmt.Sprintf("%s/%s", h.redirectHost, id))
}

func (h *Handlers) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")
	url, err := h.urlHandler.ReadURL(alias)
	if err != nil {
		h.log.Error(
			"failed to read url",
			zap.String("alias", alias),
			zap.Error(err),
		)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handlers) SaveJSONHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := models.URLRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(
			"failed to read JSON request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.urlHandler.SaveURL(req.URL)
	if err != nil {
		h.log.Error(
			"failed to save url",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, models.URLResponse{
		Result: fmt.Sprintf("%s/%s", h.redirectHost, id),
	})
}
