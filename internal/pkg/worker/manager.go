package worker

import (
	"context"
	"sync"

	"github.com/awlsring/camp/internal/pkg/logger"
)

type WorkerManagerOpts func(*WorkerManager)

type WorkerManager struct {
	workers []Worker
}

func NewWorkerManager(workers []Worker, opts ...WorkerManagerOpts) *WorkerManager {
	w := &WorkerManager{
		workers: workers,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

func (w *WorkerManager) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)

	var wg sync.WaitGroup
	for _, worker := range w.workers {
		wg.Add(1)
		go func(w Worker) {
			defer wg.Done()
			log.Debug().Msgf("starting worker %s", w.Name())
			err := w.Start(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("failed to start worker %s", w.Name())
			}
			defer w.Stop(ctx)
		}(worker)
	}
	log.Info().Msgf("started %d workers", len(w.workers))

	<-ctx.Done()
	log.Info().Msg("context done, stopping workers")
	wg.Wait()
	log.Info().Msg("all workers stopped")
	return nil
}
