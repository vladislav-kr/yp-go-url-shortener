package models

type MassURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
