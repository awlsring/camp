package rabbitmq

import (
	"context"
	"fmt"

	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/awlsring/camp/internal/pkg/worker"
	"github.com/awlsring/camp/internal/pkg/worker/rabbitmq/exchange"
	"github.com/awlsring/camp/internal/pkg/worker/rabbitmq/queue"
	amqp "github.com/rabbitmq/amqp091-go"
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

// A RabbitMQ worker that will subscribe to a queue and execute a job for each message received
type Worker struct {
	exchange *exchange.Definition
	queue    *queue.Definition
	channel  *amqp.Channel
	*worker.WorkerBase
}

// Returns a new RabbitMQ worker. The worker assumes that the exchange and queue have already been declared.
func NewWorker(channel *amqp.Channel, jobDef *JobDefinition, opts ...WorkerOpt) *Worker {
	name := fmt.Sprintf("%s-worker", jobDef.Queue.Name)

	getWork := func(ctx context.Context, jobChannel chan []byte) {
		log := logger.FromContext(ctx)
		msgs, err := channel.Consume(
			jobDef.Queue.Name, // queue
			name,              // consumer
			true,              // auto ack
			false,             // exclusive
			false,             // no local
			false,             // no wait
			nil,               // args
		)
		if err != nil {
			log.Error().Err(err).Msgf("failed to consume from queue %s", jobDef.Queue.Name)
			return
		}

		for d := range msgs {
			log.Debug().Msgf("received a message: %s", d.Body)
			jobChannel <- d.Body
		}
	}

	worker := &Worker{
		exchange:   jobDef.Exchange,
		queue:      jobDef.Queue,
		channel:    channel,
		WorkerBase: worker.NewWorkerBase(name, jobDef.ConcurrentJobs, jobDef.Job, getWork),
	}

	for _, opt := range opts {
		opt(worker)
	}

	return worker
}

func (a *Worker) Stop(ctx context.Context) error {
	return nil
}

func (w *Worker) init(ctx context.Context) error {
	log := logger.FromContext(ctx)

	log.Debug().Msg("Declaring power request exchange")
	err := w.exchange.CreateExchange(ctx, w.channel)
	if err != nil {
		log.Error().Err(err).Msg("Failed to declare power request exchange")
		return err
	}

	log.Debug().Msg("Declaring power request queue and bind")
	_, err = w.queue.CreateBindedQueue(ctx, w.channel)
	if err != nil {
		log.Error().Err(err).Msg("Failed to declare power request exchange")
		return err
	}

	log.Debug().Msg("RabbitMQ worker resources initialized")
	return nil
}

func (w *Worker) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Starting RabbitMQ worker")

	log.Debug().Msg("Initing the worker")
	err := w.init(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to init the worker")
		return err
	}

	log.Debug().Msg("Starting the base worker")
	return w.WorkerBase.Start(ctx)
}

func (w *Worker) createQueueAndBind(ctx context.Context, queue queue.Definition) (*amqp.Queue, error) {
	log := logger.FromContext(ctx)

	log.Debug().Msgf("declaring queue %s", queue.Name)
	q, err := w.channel.QueueDeclare(queue.Name, false, false, false, true, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to declare queue %s", queue.Name)
		return nil, err
	}

	log.Debug().Msgf("binding queue %s to exchange %s with routing key %s", queue.Name, queue.Exchange, queue.RoutingKey)
	err = w.channel.QueueBind(queue.Name, queue.RoutingKey, queue.Exchange, false, nil)
	if err != nil {
		log.Error().Err(err).Msgf("failed to bind queue %s to exchange %s with routing key %s", queue.Name, queue.Exchange, queue.RoutingKey)
		return nil, err
	}

	return &q, nil
}
