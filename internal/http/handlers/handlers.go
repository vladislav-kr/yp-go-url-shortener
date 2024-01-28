package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/middleware/auth"
	urlhandler "github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler"
	"go.uber.org/zap"
)

//go:generate mockery --name URLHandler
type URLHandler interface {
	ReadURL(ctx context.Context, alias string) (string, error)
	SaveURL(ctx context.Context, url string, userID string) (string, error)
	SaveURLS(ctx context.Context, urls []models.BatchRequest, userID string) ([]models.BatchResponse, error)
	Ping(ctx context.Context) error
	GetURLS(ctx context.Context, userID string) ([]models.MassURL, error)
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

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*4)
	defer cancel()

	userID := auth.UserIDFromContext(r.Context())

	id, err := h.urlHandler.SaveURL(ctx, string(data), userID)
	if err != nil {
		switch {
		case errors.Is(err, urlhandler.ErrAlreadyExists):
			h.log.Info(
				"the original url exists in the database",
				zap.String("url", string(data)),
				zap.Error(err),
			)
			render.Status(r, http.StatusConflict)
			render.PlainText(w, r, fmt.Sprintf("%s/%s", h.redirectHost, id))
			return
		default:
			h.log.Error(
				"failed to save url",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, fmt.Sprintf("%s/%s", h.redirectHost, id))
}

func (h *Handlers) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*4)
	defer cancel()

	url, err := h.urlHandler.ReadURL(ctx, alias)
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

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*4)
	defer cancel()

	userID := auth.UserIDFromContext(r.Context())

	id, err := h.urlHandler.SaveURL(ctx, req.URL, userID)
	if err != nil {
		switch {
		case errors.Is(err, urlhandler.ErrAlreadyExists):
			h.log.Info(
				"the original url exists in the database",
				zap.String("url", req.URL),
				zap.Error(err),
			)
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, models.URLResponse{
				Result: fmt.Sprintf("%s/%s", h.redirectHost, id),
			})
			return
		default:
			h.log.Error(
				"failed to save url",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, models.URLResponse{
		Result: fmt.Sprintf("%s/%s", h.redirectHost, id),
	})
}

func (h *Handlers) PingHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	if err := h.urlHandler.Ping(ctx); err != nil {
		h.log.Error("no access to database", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h *Handlers) BatchHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := []models.BatchRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(
			"failed to read JSON request body",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(req) == 0 {
		h.log.Error(
			"no data to save",
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := auth.UserIDFromContext(r.Context())

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	urls, err := h.urlHandler.SaveURLS(ctx, req, userID)
	if err != nil {
		h.log.Error(
			"failed to save url",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i := range urls {
		urls[i].ShortURL = fmt.Sprintf("%s/%s", h.redirectHost, urls[i].ShortURL)
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, urls)

}

func (h *Handlers) UserUrlsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID := auth.UserIDFromContext(r.Context())

	if len(userID) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	urls, err := h.urlHandler.GetURLS(ctx, userID)
	if err != nil {
		h.log.Error(
			"failed to read urls",
			zap.Error(err),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, urls)
}
