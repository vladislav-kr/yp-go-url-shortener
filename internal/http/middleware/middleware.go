package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Middleware struct {
	log *zap.Logger
}

func New(log *zap.Logger) *Middleware {
	return &Middleware{
		log: log,
	}
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log := m.log.With(
				zap.String("url", r.URL.Path),
				zap.String("method", r.Method),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t := time.Now()
			defer func() {
				log.Info("request completed",
					zap.Int("status", ww.Status()),
					zap.Int("bytes", ww.BytesWritten()),
					zap.Duration("duration", time.Since(t)),
				)
			}()

			next.ServeHTTP(ww, r)
		},
	)
}
