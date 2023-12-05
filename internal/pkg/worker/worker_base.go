package worker

import (
	"context"
	"sync"

	"github.com/awlsring/camp/internal/pkg/logger"
)

func SetName(name string) func(*WorkerBase) {
	return func(w *WorkerBase) {
		w.name = name
	}
}

func SetGetWorkFunc(getWork GetWorkFunc) func(*WorkerBase) {
	return func(w *WorkerBase) {
		w.getWork = getWork
	}
}

// This should be a function the will continuously get work and add it to the job channel
type GetWorkFunc func(ctx context.Context, jobChannel chan []byte)

type WorkerBase struct {
	job            Job
	name           string
	getWork        GetWorkFunc
	jobChannel     chan []byte
	concurrentJobs uint32
}

func NewWorkerBase(name string, concurrentJobs uint32, job Job, getWork GetWorkFunc) *WorkerBase {
	jobChan := make(chan []byte, concurrentJobs)
	return &WorkerBase{
		job:            job,
		jobChannel:     jobChan,
		name:           name,
		concurrentJobs: concurrentJobs,
		getWork:        getWork,
	}
}

func (w *WorkerBase) Name() string {
	return w.name
}

func (w *WorkerBase) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)

	log.Debug().Msg("Starting worker processes")
	var wg sync.WaitGroup
	for i := 0; uint32(i) < w.concurrentJobs; i++ {
		wg.Add(1)
		go w.process(ctx, &wg, i)
	}
	log.Debug().Msgf("Worker %s started", w.name)

	go func() {
		log.Debug().Msgf("Starting get work loop for %s", w.name)
		w.getWork(ctx, w.jobChannel)
	}()

	<-ctx.Done()
	log.Debug().Msgf("Worker %s context done, stopping", w.name)
	wg.Wait()
	log.Debug().Msgf("Worker %s stopped", w.name)
	return nil
}

func (w *WorkerBase) process(ctx context.Context, wg *sync.WaitGroup, procNum int) {
	log := logger.FromContext(ctx)
	defer wg.Done()

	log.Debug().Msgf("Starting worker %s process %d", w.name, procNum)
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msgf("Worker %s process %d context done, stopping", w.name, procNum)
			return
		case msg := <-w.jobChannel:
			log.Debug().Msgf("Worker %s process %d received task", w.name, procNum)

			log.Debug().Msgf("Worker %s process %d executing task", w.name, procNum)
			err := w.job.Execute(ctx, msg)
			if err != nil {
				log.Error().Err(err).Msgf("Worker %s process %d failed to execute task", w.name, procNum)
			}

			log.Debug().Msgf("Worker %s process %d completed task", w.name, procNum)
		}
	}
}
