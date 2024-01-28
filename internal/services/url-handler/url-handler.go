package urlhandler

import (
	"context"
	"errors"
	"fmt"
	netURL "net/url"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	dbkeeper "github.com/vladislav-kr/yp-go-url-shortener/internal/storages/db-keeper"
)

var (
	ErrAlreadyExists = errors.New("the value already exists")
)

//go:generate mockery --name Keeperer
type Keeperer interface {
	PostURL(ctx context.Context, url string, userID string) (string, error)
	GetURL(ctx context.Context, id string) (string, error)
	SaveURLS(ctx context.Context, urls []models.BatchRequest, userID string) ([]models.BatchResponse, error)
	GetURLS(ctx context.Context, userID string) ([]models.MassURL, error)
}

//go:generate mockery --name DBPinger
type DBPinger interface {
	PingContext(ctx context.Context) error
}

type URLHandler struct {
	storage Keeperer
	pingDB  DBPinger
}

func NewURLHandler(storage Keeperer, pingDB DBPinger) *URLHandler {
	return &URLHandler{
		storage: storage,
		pingDB:  pingDB,
	}
}

func (uh *URLHandler) ReadURL(ctx context.Context, alias string) (string, error) {

	if len(alias) == 0 {
		return "", fmt.Errorf("alias is empty")
	}

	url, err := uh.storage.GetURL(ctx, alias)
	if err != nil {
		return "", fmt.Errorf("failed to read url: %w", err)
	}

	return url, nil

}

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

func (uh *URLHandler) Ping(ctx context.Context) error {
	return uh.pingDB.PingContext(ctx)
}

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

func (uh *URLHandler) GetURLS(ctx context.Context, userID string) ([]models.MassURL, error) {
	return uh.storage.GetURLS(ctx, userID)
}
