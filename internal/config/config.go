package config

import (
	"time"
)

type Config struct {
	HTTP struct {
		Host            string
		ShutdownTimeout time.Duration
	}
	URLShortener struct {
		RedirectHost string
	}
}

// Загружает конфиг приложения
func LoadConfig() (*Config, error) {
	var cfg Config

	cfg.HTTP.ShutdownTimeout = 10 * time.Second

	return &cfg, nil
}
