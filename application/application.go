package application

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Application interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context)
}

var ShutdownTimeout = 30 * time.Second

func RunApplication(app Application) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigChan
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer shutdownCancel()

		app.Shutdown(shutdownCtx)
		wg.Done()
	}()

	wg.Add(1)

	if err := app.Run(ctx); err != nil {
		return err
	}

	wg.Wait()

	return nil
}
