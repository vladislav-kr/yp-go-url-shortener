package mapkeeper

import (
	"context"
	"fmt"
	"sync"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/lib/cryptoutils"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper/file"
)

type Keeper struct {
	mutex    sync.RWMutex
	storage  map[string]string
	filePath string
}

func New(filePath string) *Keeper {
	return &Keeper{
		storage:  map[string]string{},
		filePath: filePath,
	}
}

func (k *Keeper) PostURL(ctx context.Context, url string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		id, err := cryptoutils.GenerateRandomString(10)
		if err != nil {
			return "", err
		}

		k.mutex.Lock()
		k.storage[id] = url
		k.mutex.Unlock()
		return id, nil
	}
}

func (k *Keeper) GetURL(ctx context.Context, id string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		k.mutex.RLock()
		val, ok := k.storage[id]
		k.mutex.RUnlock()
		if !ok {
			return "", fmt.Errorf("not found")
		}

		return val, nil
	}
}
func (k *Keeper) SaveURLS(ctx context.Context, urls []models.BatchRequest) ([]models.BatchResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		batchResp := make([]models.BatchResponse, 0, len(urls))

		for _, url := range urls {
			id, err := cryptoutils.GenerateRandomString(10)
			if err != nil {
				return nil, err
			}
			k.mutex.Lock()
			k.storage[id] = url.OriginalURL
			k.mutex.Unlock()
			batchResp = append(batchResp, models.BatchResponse{
				CorrelationId: url.CorrelationId,
				ShortURL:      id,
			})
		}

		return batchResp, nil
	}
}

func (k *Keeper) LoadFromFile() error {
	if len(k.filePath) == 0 {
		return nil
	}

	k.mutex.Lock()
	defer k.mutex.Unlock()

	c, err := file.NewConsumer(k.filePath)
	if err != nil {
		return err
	}

	var url *models.FileURL
	for c.More() {
		url = &models.FileURL{}
		c.Decode(url)
		k.storage[url.ShortURL] = url.OriginalURL
	}

	return c.Close()
}

func (k *Keeper) SaveToFile() error {
	if len(k.filePath) == 0 {
		return nil
	}

	k.mutex.Lock()
	defer k.mutex.Unlock()

	if len(k.storage) < 1 {
		return nil
	}

	p, err := file.NewProducer(k.filePath)
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
