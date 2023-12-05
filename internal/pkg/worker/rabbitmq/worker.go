package rabbitmq

import (
	"context"
	"fmt"
	"sync"

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
		w.name = name
	}
}

type Worker struct {
	jobChannel     chan *amqp.Delivery
	exchange       *exchange.Definition
	queue          *queue.Definition
	concurrentJobs uint32
	runningJobs    uint32
	job            worker.Job
	name           string
	channel        *amqp.Channel
}

func NewWorker(channel *amqp.Channel, jobDef *JobDefinition, opts ...WorkerOpt) *Worker {
	jobChan := make(chan *amqp.Delivery, jobDef.ConcurrentJobs)
	worker := &Worker{
		jobChannel:     jobChan,
		exchange:       jobDef.Exchange,
		queue:          jobDef.Queue,
		concurrentJobs: jobDef.ConcurrentJobs,
		job:            jobDef.Job,
		name:           fmt.Sprintf("%s-worker", jobDef.Queue.Name),
		channel:        channel,
	}

	for _, opt := range opts {
		opt(worker)
	}

	return worker
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
		case job := <-a.jobChannel:
			log.Debug().Msgf("worker %s process %d received task", a.name, procNum)

			log.Debug().Msgf("worker %s process %d executing task", a.name, procNum)
			err := a.job.Execute(ctx, job.Body)
			if err != nil {
				log.Error().Err(err).Msgf("worker %s process %d failed to execute task", a.name, procNum)
			}

			log.Debug().Msgf("worker %s process %d completed task", a.name, procNum)
		}
	}
}

func (a *Worker) Stop(ctx context.Context) error {
	return nil
}

func (a *Worker) Name() string {
	return a.name
}

func (w *Worker) init(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("initializing rabbitmq worker")

	log.Debug().Msg("declaring power request exchange")
	err := w.exchange.CreateExchange(ctx, w.channel)
	if err != nil {
		log.Error().Err(err).Msg("failed to declare power request exchange")
		return err
	}

	_, err = w.queue.CreateBindedQueue(ctx, w.channel)
	if err != nil {
		log.Error().Err(err).Msg("failed to declare power request exchange")
		return err
	}

	return nil
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

func (w *Worker) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)

	err := w.init(ctx)
	if err != nil {
		return err
	}

	log.Debug().Msg("starting rabbitmq act")
	var wg sync.WaitGroup
	for i := 0; uint32(i) < w.concurrentJobs; i++ {
		wg.Add(1)
		go w.process(ctx, &wg, i)
	}
	log.Debug().Msg("rabbitmq act started")

	msgs, err := w.channel.Consume(
		w.queue.Name, // queue
		w.name,       // consumer
		true,         // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	// TODO: handle max jobs
	go func() {
		for d := range msgs {
			log.Debug().Msgf("received a message: %s", d.Body)
			w.jobChannel <- &d
			// d.Ack(false)
		}
	}()

	<-ctx.Done()
	log.Debug().Msg("rabbitmq act context done, stopping")
	wg.Wait()
	log.Debug().Msg("rabbitmq act stopped")
	return nil

}
