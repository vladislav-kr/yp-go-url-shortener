package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/app"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/config"
	"golang.org/x/sync/errgroup"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fail load config: %w", err))
	}

	parseFlags(&cfg.HTTP.Host, &cfg.URLShortener.RedirectHost)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCtx, sigCancel := signal.NotifyContext(ctx, os.Interrupt)
	defer sigCancel()

	errGr, errGrCtx := errgroup.WithContext(sigCtx)

	srv := app.NewServer(
		app.Option{
			Host:         cfg.HTTP.Host,
			RedirectHost: cfg.URLShortener.RedirectHost,
		},
	)

	errGr.Go(func() error {
		return srv.Run()
	})

	errGr.Go(func() error {
		<-errGrCtx.Done()
		return srv.Stop()
	})

	if err := errGr.Wait(); err != nil {
		log.Fatal(err)
	}

}
