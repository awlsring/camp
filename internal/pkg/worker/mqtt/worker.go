package mqtt

import (
	"context"
	"fmt"

	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/awlsring/camp/internal/pkg/worker"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _ worker.Worker = &Worker{}

type WorkerOpt func(*Worker)

func WithName(name string) WorkerOpt {
	return func(w *Worker) {
		worker.SetName(name)(w.WorkerBase)
	}
}

func WithGetWorkFunc(getWork worker.GetWorkFunc) WorkerOpt {
	return func(w *Worker) {
		worker.SetGetWorkFunc(getWork)(w.WorkerBase)
	}
}

// An MQTT worker that will subscribe to a topic and execute a job for each message received
type Worker struct {
	topic  string
	client mqtt.Client
	*worker.WorkerBase
}

func NewWorker(client mqtt.Client, jobDef *JobDefinition, opts ...WorkerOpt) *Worker {
	name := fmt.Sprintf("%s-worker", jobDef.Topic)
	getWork := func(ctx context.Context, jobChannel chan []byte) {
		log := logger.FromContext(ctx)
		client.Subscribe(jobDef.Topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			log.Debug().Msgf("received message: %s", msg.Payload())
			jobChannel <- msg.Payload()
		})
	}
	w := &Worker{
		topic:      jobDef.Topic,
		client:     client,
		WorkerBase: worker.NewWorkerBase(name, jobDef.ConcurrentJobs, jobDef.Job, getWork),
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

func (a *Worker) Stop(ctx context.Context) error {
	return nil
}
