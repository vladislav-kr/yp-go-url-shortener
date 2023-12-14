package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlerer interface {
	SaveHandler(w http.ResponseWriter, r *http.Request)
	RedirectHandler(w http.ResponseWriter, r *http.Request)
}

// Конфигурирует главный роутер
func NewRouter(
	handlers Handlerer,
) *chi.Mux {

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
		r.Post("/", handlers.SaveHandler)
		r.Get("/{id}", handlers.RedirectHandler)
	})

	return router
}
