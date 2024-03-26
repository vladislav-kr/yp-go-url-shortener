package models

// FileURL структура хранения в файле.
type FileURL struct {
	ShortURL    string `json:"shortUrl"`
	OriginalURL string `json:"originalUrl"`
}
