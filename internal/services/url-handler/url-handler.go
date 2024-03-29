// urlhandler отвечает за реализацию бизнес логики
package urlhandler

import (
	"context"
	"errors"
	"fmt"
	netURL "net/url"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/services/url-handler/deleter"
	dbkeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/db-keeper"
)

// Ошибки сервисного слоя
var (
	ErrAlreadyExists = errors.New("the value already exists")
	ErrURLRemoved    = errors.New("url has already been deleted")
)

// Keeperer интерфейс хранилища данных.
//
//go:generate mockery --name Keeperer
type Keeperer interface {
	PostURL(ctx context.Context, url string, userID string) (string, error)
	GetURL(ctx context.Context, id string) (string, error)
	SaveURLS(ctx context.Context, urls []models.BatchRequest, userID string) ([]models.BatchResponse, error)
	GetURLS(ctx context.Context, userID string) ([]models.MassURL, error)
	DeleteURLS(ctx context.Context, shortURLS []models.DeleteURL)
}

// DBPinger интерфейс проверки доступности хранилища.
//
//go:generate mockery --name DBPinger
type DBPinger interface {
	PingContext(ctx context.Context) error
}

// URLHandler хранит объекты, необходимые для реализации бизнес логики
type URLHandler struct {
	storage Keeperer
	pingDB  DBPinger
	deleter *deleter.Deleter
}

// NewURLHandler конструктор URLHandler.
func NewURLHandler(storage Keeperer, pingDB DBPinger, deleter *deleter.Deleter) *URLHandler {
	return &URLHandler{
		storage: storage,
		pingDB:  pingDB,
		deleter: deleter,
	}
}

// ReadURL чтение оригинального URL.
func (uh *URLHandler) ReadURL(ctx context.Context, alias string) (string, error) {

	if len(alias) == 0 {
		return "", fmt.Errorf("alias is empty")
	}

	url, err := uh.storage.GetURL(ctx, alias)
	if err != nil {
		if errors.Is(err, dbkeeper.ErrURLRemoved) {
			return "", ErrURLRemoved
		}
		return "", fmt.Errorf("failed to read url: %w", err)
	}

	return url, nil

}

// SaveURL сохранение сокращенного URL.
func (uh *URLHandler) SaveURL(ctx context.Context, url string, userID string) (string, error) {

	alias := ""
	if _, err := netURL.ParseRequestURI(url); err != nil {
		return alias, fmt.Errorf("invalid url: %w", err)
	}

	alias, err := uh.storage.PostURL(ctx, url, userID)
	if err != nil {
		switch {
		case errors.Is(err, dbkeeper.ErrAlreadyExists):
			return alias, ErrAlreadyExists
		default:
			return alias, fmt.Errorf("failed to save url: %w", err)
		}
	}

	return alias, nil
}

// Ping проверка доступности хранилища.
func (uh *URLHandler) Ping(ctx context.Context) error {
	return uh.pingDB.PingContext(ctx)
}

// SaveURLS массовое сохранение URL.
func (uh *URLHandler) SaveURLS(
	ctx context.Context,
	urls []models.BatchRequest,
	userID string,
) (
	[]models.BatchResponse,
	error,
) {
	return uh.storage.SaveURLS(ctx, urls, userID)
}

// GetURLS список сокращенных URL пользователя.
func (uh *URLHandler) GetURLS(ctx context.Context, userID string) ([]models.MassURL, error) {
	return uh.storage.GetURLS(ctx, userID)
}

// DeleteURLS удаление URL.
func (uh *URLHandler) DeleteURLS(ctx context.Context, shortURLS []string, userID string) {
	select {
	case <-ctx.Done():
		return
	default:
		uh.deleter.AddMessages(shortURLS, userID)
	}
}
