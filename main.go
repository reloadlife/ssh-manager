package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mamad.dev/ssh-manager/manager"
	"os"
	"os/signal"
	"syscall"
)

func run(ctx context.Context) error {
	return manager.Application(ctx)
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := run(ctx); err != nil {
			log.Error(err)
		}
		cancel()
	}()

	select {
	case <-sig:
		cancel()
	case <-ctx.Done():
	}
}
