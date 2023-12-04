package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	HTTP struct {
		Host            string        `env:"SERVER_ADDRESS"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
		ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"4s"`
		WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"4s"`
		IdleTimeout     time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"15s"`
	}
	URLShortener struct {
		RedirectHost string `env:"BASE_URL"`
	}
}

// Загружает конфиг приложения
func LoadConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("fail read env: %w", err)
	}

	return &cfg, nil
}
