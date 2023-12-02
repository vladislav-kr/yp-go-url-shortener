package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/handlers/redirect"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/handlers/save"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/keeper"
)

// Конфигурирует главный роутер
func NewRouter() (*chi.Mux, error) {
	stor := keeper.New()
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.URLFormat,
	)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	router.Route("/", func(r chi.Router) {
		r.Post("/", save.NewSaveHandler(stor))
		r.Get("/{id}", redirect.NewRedirectHandler(stor))
	})

	return router, nil
}
