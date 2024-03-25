// mapkeeper in-memory хранилище
package mapkeeper

import (
	"context"
	"fmt"
	"sync"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/lib/cryptoutils"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/storages/map-keeper/file"
)

// Keeper хранит данные для in-memory хранилища.
type Keeper struct {
	mutex    sync.RWMutex
	storage  map[string]string
	filePath string
}

// New конструктор Keeper.
func New(filePath string) *Keeper {
	return &Keeper{
		storage:  map[string]string{},
		filePath: filePath,
	}
}

// DeleteURLS удаление URL
func (k *Keeper) DeleteURLS(_ context.Context, _ []models.DeleteURL) {

}

// PostURL сохранение сокращенного URL.
func (k *Keeper) PostURL(ctx context.Context, url string, _ string) (string, error) {
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

// GetURL чтение оригинального URL.
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

// SaveURLS массовое сохранение URL.
func (k *Keeper) SaveURLS(ctx context.Context, urls []models.BatchRequest, _ string) ([]models.BatchResponse, error) {
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
				CorrelationID: url.CorrelationID,
				ShortURL:      id,
			})
		}

		return batchResp, nil
	}
}

// LoadFromFile загружает данные из файла.
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

// SaveToFile сохраняет данные в файл .
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

func (k *Keeper) GetURLS(_ context.Context, _ string) ([]models.MassURL, error) {
	return nil, fmt.Errorf("method not implemented")
}
