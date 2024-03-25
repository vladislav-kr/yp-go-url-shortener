package main

import (
	"context"
	"log"

	"go.uber.org/zap"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/app"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/config"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/logger"

)

func main() {

	// Загружает конфиг из env
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("fail load config: %v", err)
	}

	log := logger.MustLogger(cfg.App.LogLevel)
	defer log.Sync()
	log.Info("launching a url shortener...")
	log.Debug("debug messages enabled")

	// Дополним конфиг из флагов, если env переменные не заданы
	parseFlags(
		&cfg.HTTP.Host,
		&cfg.URLShortener.RedirectHost,
		&cfg.Storage.File.PATH,
		&cfg.Storage.Postgres.DNS,
	)

	// Основной контекст api сервера
	// Не отменяется при отмене errgroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	urlShortener, err := app.NewURLShortener(
		ctx,
		log,
		app.Option{
			Host:            cfg.HTTP.Host,
			RedirectHost:    cfg.URLShortener.RedirectHost,
			ReadTimeout:     cfg.HTTP.ReadTimeout,
			WriteTimeout:    cfg.HTTP.WriteTimeout,
			IdleTimeout:     cfg.HTTP.IdleTimeout,
			ShutdownTimeout: cfg.HTTP.ShutdownTimeout,
			StorageFilePath: cfg.Storage.File.PATH,
			StorageDBDNS:    cfg.Storage.Postgres.DNS,
		},
	)
	if err != nil {
		log.Error("failed to configure server", zap.Error(err))
		return
	}

	if err := urlShortener.Run(ctx); err != nil {
		log.Error(err.Error())
	}

}
