package keeper

import (
	"fmt"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/lib/cryptoutils"
)

type Keeperer interface {
	PostURL(url string) (string, error)
	GetURL(id string) (string, error)
}

type Keeper struct {
	storage map[string]string
}

func New() *Keeper{
	return &Keeper{
		storage: map[string]string{},
	}
}

func (k *Keeper) PostURL(url string) (string, error) {

	id, err := cryptoutils.GenerateRandomString(10)
	if err != nil {

		return "", err
	}

	k.storage[id] = url

	return id, nil
}

func (k *Keeper) GetURL(id string) (string, error) {

	val, ok := k.storage[id]
	if !ok {
		return "", fmt.Errorf("not found")
	}

	return val, nil
}
