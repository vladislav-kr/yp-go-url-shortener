package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/handlers"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/http/middleware"
	urlHandler "github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler"
	mapKeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper"
	"go.uber.org/zap/zaptest"
)

func TestNewRouter(t *testing.T) {
	log := zaptest.NewLogger(t)
	h := handlers.NewHandlers(
		log,
		urlHandler.NewURLHandler(
			mapKeeper.New(),
		),
		"http://localhost:8080",
	)
	m := middleware.New(log)
	r := NewRouter(h, m)

	assert.NotEmpty(t, r)
}
