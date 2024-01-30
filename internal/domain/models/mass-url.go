package models

type MassURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type MassDeleteURL struct {
	ShortURLS []string `json:"short_urls"`
	UserID    string   `json:"user_id"`
}

type DeleteURL struct {
	ShortURL string `json:"short_url"`
	UserID   string `json:"user_id"`
}
