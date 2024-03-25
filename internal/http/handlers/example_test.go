package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	urlhandler "github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler/deleter"
	mapkeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper"
)

func ExampleHandlers_SaveHandler() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем in-memory хранилище
	storage := mapkeeper.New("")

	// Создаем обработчик сервисного слоя
	urlHandler := urlhandler.NewURLHandler(
		storage,
		nil,
		deleter.NewDeleter(ctx, 10, func(urls []models.DeleteURL) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			storage.DeleteURLS(ctx, urls)
		}),
	)

	// Создаем инстанцию http обработчиков
	h := NewHandlers(
		zap.L(),
		urlHandler,
		"http://localhost:8080",
	)

	router := chi.NewRouter()
	// Регистрируем обработчик на роутере
	router.Post("/", h.SaveHandler)

	ctxReq, cancelReq := context.WithTimeout(ctx, time.Second*2)
	defer cancelReq()

	// Создаем запрос на сокращение URL
	req, _ := http.NewRequestWithContext(
		ctxReq,
		http.MethodPost,
		"/",
		strings.NewReader("https://go.dev/"),
	)
	ww := httptest.NewRecorder()

	router.ServeHTTP(ww, req)
	ww.Result().Body.Close()
	fmt.Println(ww.Result().StatusCode)
	// Output:
	// 201

}
