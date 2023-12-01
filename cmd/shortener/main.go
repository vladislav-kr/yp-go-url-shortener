package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/app"
	"golang.org/x/sync/errgroup"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCtx, sigCancel := signal.NotifyContext(ctx, os.Interrupt)
	defer sigCancel()

	errGr, errGrCtx := errgroup.WithContext(sigCtx)

	srv := app.NewServer(":8080")

	errGr.Go(func() error {
		return srv.Run()
	})

	errGr.Go(func() error {
		<-errGrCtx.Done()
		return srv.Stop()
	})

	if err := errGr.Wait(); err != nil {
		panic(err)
	}

}
