package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/middleware/auth"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/middleware/compress"
	"go.uber.org/zap"
)

type Middleware struct {
	log  *zap.Logger
	auth *auth.Auth
}

func New(log *zap.Logger, auth *auth.Auth) *Middleware {
	return &Middleware{
		log:  log,
		auth: auth,
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

func (m *Middleware) NewCompressHandler(contentTypes []string) func(next http.Handler) http.Handler {
	compressPool, err := compress.NewCompressPool(contentTypes)
	if err != nil {
		m.log.Error("failed to create compressor pool", zap.Error(err))
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if compressPool == nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
			// который будем передавать следующей функции
			ow := w

			// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
				cw := compressPool.NewCompressWriter(w)
				// меняем оригинальный http.ResponseWriter на новый
				ow = cw
				// не забываем отправить клиенту все сжатые данные после завершения middleware
				defer cw.Close()
			}

			// проверяем, что клиент отправил серверу сжатые данные в формате gzip
			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
				cr := compressPool.NewCompressReader(r.Body)

				r.Body = cr
				defer cr.Close()
			}

			// передаём управление хендлеру
			next.ServeHTTP(ow, r)
		})
	}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			cookie, err := auth.CookieFromRequest(r)
			if err != nil {
				userID := uuid.New().String()
				cookie, err := m.auth.CreateCookie(time.Hour*3, userID)
				if err != nil {
					m.log.Error("failed to create new cookie", zap.Error(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, cookie)
				next.ServeHTTP(w, r)
				return
			}

			claims, err := m.auth.Validate(cookie.Value)
			if err != nil {
				userID := uuid.New().String()
				m.log.Error("token is invalid", zap.Error(err))
				cookie, err := m.auth.CreateCookie(time.Hour*3, userID)
				if err != nil {
					m.log.Error("failed to create new cookie", zap.Error(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, cookie)
				next.ServeHTTP(w, r)
				return
			}

			ctx := auth.ContextWithUserID(r.Context(), claims.UserID)
			rr := r.WithContext(ctx)

			next.ServeHTTP(w, rr)
		},
	)
}
