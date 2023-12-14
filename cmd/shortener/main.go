package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/app"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/config"
	"golang.org/x/sync/errgroup"
)

func main() {

	// Загружает конфиг из env
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("fail load config: %v", err)
	}

	// Дополним конфиг из флагов, если env переменные не заданы
	parseFlags(&cfg.HTTP.Host, &cfg.URLShortener.RedirectHost)

	// Основной контекст api сервера
	// Не отменяется при отмене errgroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Контекст прослушивающий сигналы прерывания OS
	sigCtx, sigCancel := signal.NotifyContext(ctx,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer sigCancel()

	// группа для запуска и остановки сервера по сигналу
	errGr, errGrCtx := errgroup.WithContext(sigCtx)

	srv := app.NewServer(
		app.Option{
			Host:         cfg.HTTP.Host,
			RedirectHost: cfg.URLShortener.RedirectHost,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
			IdleTimeout:  cfg.HTTP.IdleTimeout,
		},
	)

	errGr.Go(func() error {
		return srv.Run()
	})

	errGr.Go(func() error {
		<-errGrCtx.Done()

		ctx, cancel := context.WithTimeout(
			context.Background(),
			cfg.HTTP.ShutdownTimeout,
		)
		defer cancel()

		return srv.Stop(ctx)
	})

	if err := errGr.Wait(); err != nil {
		log.Fatal(err)
	}

}
