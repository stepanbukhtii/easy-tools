package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ShutdownTimeout = 30 * time.Second

type Application struct {
	sigChan chan os.Signal
}

func NewApplication() Application {
	return Application{
		sigChan: make(chan os.Signal, 1),
	}
}

func (a *Application) Run(shutdownFunc func(context.Context)) context.Context {
	signal.Notify(a.sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-a.sigChan
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer shutdownCancel()

		shutdownFunc(shutdownCtx)
	}()

	return ctx
}

func RunApplication(shutdownFunc func(context.Context)) context.Context {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigChan
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer shutdownCancel()

		shutdownFunc(shutdownCtx)
	}()

	return ctx
}
