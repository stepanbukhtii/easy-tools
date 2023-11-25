package worker

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type WorkerExample struct {
	name string
}

func (w *WorkerExample) Name() string {
	return w.name
}

func (w *WorkerExample) Run(ctx context.Context) {
	fmt.Println(w.name, "start")

	<-ctx.Done()
	fmt.Println(w.name, "finish")
}

type WorkerPanic struct {
	name string
}

func (w *WorkerPanic) Name() string {
	return w.name
}

func (w *WorkerPanic) Run(ctx context.Context) {
	fmt.Println(w.name, "start")

	time.Sleep(400 * time.Millisecond)

	var foo map[string]string
	foo["key"] = "value"

	<-ctx.Done()
	fmt.Println(w.name, "finish")
}

func TestWorkerHub(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	worHub := NewHub()
	worHub.AddWorker("test1", &WorkerExample{name: "test1"})
	worHub.AddWorker("test2", &WorkerExample{name: "test2"})
	worHub.AddWorker("test3", &WorkerPanic{name: "test3"})
	worHub.Run(ctx)
}
