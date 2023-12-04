package rabbitmq_worker

import (
	"context"
	"sync"

	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/awlsring/camp/internal/pkg/rabbitmq_worker/job"
	amqp "github.com/rabbitmq/amqp091-go"
)

type WorkerManagerOpts func(*WorkerManager)

type WorkerManager struct {
	workers []*Worker
}

func NewWorkerManager(channel *amqp.Channel, jobDefs []*job.Definition, opts ...WorkerManagerOpts) *WorkerManager {
	w := &WorkerManager{}

	for _, opt := range opts {
		opt(w)
	}

	for _, jobDef := range jobDefs {
		w.workers = append(w.workers, NewWorker(channel, jobDef.Queue, jobDef.Exchange, jobDef.ConcurrentJobs, jobDef.Job))
	}

	return w
}

func (w *WorkerManager) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)

	var wg sync.WaitGroup
	for _, worker := range w.workers {
		wg.Add(1)
		go func(a *Worker) {
			defer wg.Done()
			err := a.Start(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to start worker")
			}
		}(worker)
	}
	log.Info().Msgf("started %d workers", len(w.workers))

	<-ctx.Done()
	log.Info().Msg("context done, stopping workers")
	wg.Wait()
	log.Info().Msg("all workers stopped")
	return nil
}
