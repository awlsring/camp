package mqtt

import (
	"context"
	"fmt"
	"sync"

	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/awlsring/camp/internal/pkg/worker"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _ worker.Worker = &Worker{}

type WorkerOpt func(*Worker)

func WithName(name string) WorkerOpt {
	return func(w *Worker) {
		w.name = name
	}
}

type Worker struct {
	jobChannel     chan mqtt.Message
	concurrentJobs uint32
	runningJobs    uint32
	topic          string
	name           string
	job            worker.Job
	client         mqtt.Client
}

func NewWorker(client mqtt.Client, jobDef *JobDefinition, opts ...WorkerOpt) *Worker {
	jobChan := make(chan mqtt.Message, jobDef.ConcurrentJobs)
	worker := &Worker{
		jobChannel:     jobChan,
		concurrentJobs: jobDef.ConcurrentJobs,
		name:           fmt.Sprintf("%s-worker", jobDef.Topic),
		topic:          jobDef.Topic,
		client:         client,
	}

	for _, opt := range opts {
		opt(worker)
	}

	return worker
}

func (a *Worker) Name() string {
	return a.name
}

func (a *Worker) Stop(ctx context.Context) error {
	return nil
}

func (a *Worker) process(ctx context.Context, wg *sync.WaitGroup, procNum int) {
	log := logger.FromContext(ctx)

	defer wg.Done()

	log.Debug().Msgf("starting worker %s process %d", a.name, procNum)
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msgf("worker %s process %d context done, stopping", a.name, procNum)
			return
		case msg := <-a.jobChannel:
			log.Debug().Msgf("worker %s process %d received task", a.name, procNum)

			log.Debug().Msgf("worker %s process %d executing task", a.name, procNum)
			err := a.job.Execute(ctx, msg.Payload())
			if err != nil {
				log.Error().Err(err).Msgf("worker %s process %d failed to execute task", a.name, procNum)
			}

			log.Debug().Msgf("worker %s process %d completed task", a.name, procNum)
		}
	}
}

func (w *Worker) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)

	log.Debug().Msg("Starting mqtt worker")
	var wg sync.WaitGroup
	for i := 0; uint32(i) < w.concurrentJobs; i++ {
		wg.Add(1)
		go w.process(ctx, &wg, i)
	}
	log.Debug().Msg("MQTT worker started")

	go func() {
		w.client.Subscribe(w.topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			log.Debug().Msgf("received message: %s", msg.Payload())
			w.jobChannel <- msg
		})
	}()

	<-ctx.Done()
	log.Debug().Msg("MQTT worker context done, stopping")
	wg.Wait()
	log.Debug().Msg("MQTT worker stopped")
	return nil

}
