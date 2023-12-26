package mapkeeper

import (
	"fmt"
	"sync"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/lib/cryptoutils"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper/file"
)

type Keeper struct {
	mutex   sync.RWMutex
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

	k.mutex.Lock()
	k.storage[id] = url
	k.mutex.Unlock()
	return id, nil
}

func (k *Keeper) GetURL(id string) (string, error) {
	k.mutex.RLock()
	val, ok := k.storage[id]
	k.mutex.RUnlock()
	if !ok {
		return "", fmt.Errorf("not found")
	}

	return val, nil
}

func (k *Keeper) LoadFromFile(path string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	c, err := file.NewConsumer(path)
	if err != nil {
		return err
	}

	url := &models.FileURL{}

	for c.More() {
		c.Decode(url)
		k.storage[url.ShortURL] = url.OriginalURL
	}

	return c.Close()
}

func (k *Keeper) SaveToFile(path string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	if len(k.storage) < 1 {
		return nil
	}

	p, err := file.NewProducer(path)
	if err != nil {
		return err
	}

	url := &models.FileURL{}
	for shortURL, originURL := range k.storage {
		url.ShortURL = shortURL
		url.OriginalURL = originURL
		p.Write(url)
	}

	return p.Close()
}
