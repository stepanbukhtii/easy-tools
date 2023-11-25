package workerhub

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"sync"
)

type Worker interface {
	Run(ctx context.Context)
	Name() string
}

type Initializer interface {
	Init(ctx context.Context) error
}

type WorkerHub map[string]Worker

func NewWorkerHub() WorkerHub {
	return make(WorkerHub)
}

func (wh WorkerHub) AddWorker(name string, worker Worker) {
	wh[name] = worker
}

func (wh WorkerHub) Init(ctx context.Context) error {
	for _, worker := range wh {
		workerInit, ok := worker.(Initializer)
		if !ok {
			continue
		}
		if err := workerInit.Init(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (wh WorkerHub) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, worker := range wh {
		go wh.RunWorker(ctx, &wg, worker)
	}
	wg.Wait()
}

func (wh WorkerHub) RunWorker(ctx context.Context, wg *sync.WaitGroup, worker Worker) {
	defer func() {
		wg.Done()

		r := recover()
		if r == nil {
			return
		}

		log.Error().
			Err(fmt.Errorf("%v", r)).
			Str("worker_name", worker.Name()).
			Msg("panic recovered in worker")

		wh.RunWorker(ctx, wg, worker)
	}()

	wg.Add(1)
	worker.Run(ctx)
}
