package urlhandler

import (
	"context"
	"fmt"
	netURL "net/url"
)

//go:generate mockery --name Keeperer
type Keeperer interface {
	PostURL(url string) (string, error)
	GetURL(id string) (string, error)
}

//go:generate mockery --name DBKeeperer
type DBKeeperer interface {
	Ping(ctx context.Context) error
}

type URLHandler struct {
	storage   Keeperer
	dbStorage DBKeeperer
}

func NewURLHandler(storage Keeperer, dbStorage DBKeeperer) *URLHandler {
	return &URLHandler{
		storage:   storage,
		dbStorage: dbStorage,
	}
}

func (uh *URLHandler) ReadURL(alias string) (string, error) {

	if len(alias) == 0 {
		return "", fmt.Errorf("alias is empty")
	}

	url, err := uh.storage.GetURL(alias)
	if err != nil {
		return "", fmt.Errorf("failed to read url: %w", err)
	}

	return url, nil

}

func (uh *URLHandler) SaveURL(url string) (string, error) {

	alias := ""
	if _, err := netURL.ParseRequestURI(url); err != nil {
		return alias, fmt.Errorf("invalid url: %w", err)
	}

	alias, err := uh.storage.PostURL(url)
	if err != nil {
		return alias, fmt.Errorf("failed to save url: %w", err)
	}
	return alias, nil
}

func (uh *URLHandler) Ping(ctx context.Context) error {
	return uh.dbStorage.Ping(ctx)
}
