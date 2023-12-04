package rabbitmq_worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/awlsring/camp/internal/pkg/rabbitmq_worker/exchange"
	"github.com/awlsring/camp/internal/pkg/rabbitmq_worker/job"
	"github.com/awlsring/camp/internal/pkg/rabbitmq_worker/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	jobChannel     chan *amqp.Delivery
	exchange       *exchange.Definition
	queue          *queue.Definition
	concurrentJobs uint32
	runningJobs    uint32
	job            job.Job
	name           string
	channel        *amqp.Channel
}

func NewWorker(channel *amqp.Channel, queue *queue.Definition, exchange *exchange.Definition, concurrentJobs uint32, job job.Job) *Worker {
	jobChan := make(chan *amqp.Delivery, concurrentJobs)
	return &Worker{
		jobChannel:     jobChan,
		exchange:       exchange,
		queue:          queue,
		concurrentJobs: concurrentJobs,
		job:            job,
		name:           fmt.Sprintf("%s-worker", queue.Name),
		channel:        channel,
	}
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

	go func() {
		for d := range msgs {
			log.Debug().Msgf("received a message: %s", d.Body)
			w.jobChannel <- &d
		}
	}()

	<-ctx.Done()
	log.Debug().Msg("rabbitmq act context done, stopping")
	wg.Wait()
	log.Debug().Msg("rabbitmq act stopped")
	return nil

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Debug().Msg("rabbitmq act context done, stopping")
	// 		err := a.Stop(ctx)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}

	// 	// wait for tasks to come on queue and send to chan
	// 	runningJobs := a.runningJobs
	// 	log.Debug().Msgf("Current running jobs: %d", runningJobs)

	// 	if runningJobs >= a.concurrentJobs {
	// 		log.Debug().Msg("rabbitmq act running jobs at capacity, waiting for jobs to complete")
	// 		time.Sleep(1 * time.Second)
	// 		continue
	// 	}

	// 	//poll
	// 	// ????

	// }
}
