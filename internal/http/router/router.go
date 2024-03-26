package router

import (
	"net/http"
	"net/http/pprof"

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

	// Регистрация pprof-обработчиков
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	router.Handle("/debug/pprof/block", pprof.Handler("block"))
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	return router
}
