package mapkeeper

import (
	"fmt"
	"sync"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/lib/cryptoutils"
)

type Keeper struct {
	m       sync.RWMutex
	storage map[string]string
}

func New() *Keeper {
	return &Keeper{
		storage: map[string]string{},
	}
}

func (k *Keeper) PostURL(url string) (string, error) {

	id, err := cryptoutils.GenerateRandomString(10)
	if err != nil {
		return "", err
	}

	k.m.Lock()
	k.storage[id] = url
	k.m.Unlock()
	return id, nil
}

func (k *Keeper) GetURL(id string) (string, error) {
	k.m.RLock()
	val, ok := k.storage[id]
	k.m.RUnlock()
	if !ok {
		return "", fmt.Errorf("not found")
	}

	return val, nil
}
