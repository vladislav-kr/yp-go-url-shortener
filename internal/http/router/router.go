package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/handlers"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/middleware"
)

// Конфигурирует главный роутер
func NewRouter(
	h *handlers.Handlers,
	m *middleware.Middleware,
) *chi.Mux {

	router := chi.NewRouter()

	router.Use(
		chiMiddleware.Recoverer,
		chiMiddleware.URLFormat,
		m.NewCompressHandler([]string{
			"application/json",
			"text/html",
		}),
		m.Auth,
		m.Logger,
	)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	router.Get("/ping", h.PingHandler)
	router.Get("/{id}", h.RedirectHandler)
	router.Post("/", h.SaveHandler)
	router.Post("/api/shorten", h.SaveJSONHandler)
	router.Post("/api/shorten/batch", h.BatchHandler)
	router.Get("/api/user/urls", h.UserUrlsHandler)
	router.Delete("/api/user/urls", h.DeleteURLS)

	return router
}
