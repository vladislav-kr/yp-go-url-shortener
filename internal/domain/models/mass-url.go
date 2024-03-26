package models

// MassURL массовое удаление сокращенных URL.
type MassURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// MassDeleteURL группа URL для удаления по UserID.
type MassDeleteURL struct {
	ShortURLS []string `json:"short_urls"`
	UserID    string   `json:"user_id"`
}

// DeleteURL НН URL для удаления.
type DeleteURL struct {
	ShortURL string `json:"short_url"`
	UserID   string `json:"user_id"`
}
