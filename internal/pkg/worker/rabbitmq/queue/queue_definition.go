package queue

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueOpt func(*Definition)

type Definition struct {
	Name        string
	RoutingKey  string
	Exchange    string
	Durable     bool
	AutoDelete  bool
	Exclusive   bool
	NoWait      bool
	QueueArgs   amqp.Table
	BindingArgs amqp.Table
}

func (q Definition) CreateBindedQueue(ctx context.Context, channel *amqp.Channel) (*amqp.Queue, error) {
	log := logger.FromContext(ctx)

	log.Debug().Msgf("declaring queue %s", q.Name)
	queue, err := channel.QueueDeclare(q.Name, q.Durable, q.AutoDelete, q.Exclusive, q.NoWait, q.QueueArgs)
	if err != nil {
		log.Error().Err(err).Msgf("failed to declare queue %s", queue.Name)
		return nil, err
	}
	log.Debug().Msgf("declared queue %s", queue.Name)

	log.Debug().Msgf("binding queue %s to exchange %s with routing key %s", q.Name, q.Exchange, q.RoutingKey)
	err = channel.QueueBind(q.Name, q.RoutingKey, q.Exchange, q.NoWait, q.BindingArgs)
	if err != nil {
		log.Error().Err(err).Msgf("failed to bind queue %s to exchange %s with routing key %s", q.Name, q.Exchange, q.RoutingKey)
		return nil, err
	}
	log.Debug().Msgf("bound queue %s to exchange %s with routing key %s", q.Name, q.Exchange, q.RoutingKey)

	return &queue, nil
}

func WithDurable(durable bool) QueueOpt {
	return func(e *Definition) {
		e.Durable = durable
	}
}

func WithAutoDelete(autoDelete bool) QueueOpt {
	return func(e *Definition) {
		e.AutoDelete = autoDelete
	}
}

func WithExclusive(ex bool) QueueOpt {
	return func(e *Definition) {
		e.Exclusive = ex
	}
}

func WithNoWait(noWait bool) QueueOpt {
	return func(e *Definition) {
		e.NoWait = noWait
	}
}

func WithQueueArgs(args amqp.Table) QueueOpt {
	return func(e *Definition) {
		e.QueueArgs = args
	}
}

func WithBindingArgs(args amqp.Table) QueueOpt {
	return func(e *Definition) {
		e.BindingArgs = args
	}
}

func NewDefinition(name, routingKey, exchange string, opts ...QueueOpt) *Definition {
	def := &Definition{
		Name:       name,
		RoutingKey: routingKey,
		Exchange:   exchange,
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     true,
		QueueArgs:  nil,
	}

	for _, opt := range opts {
		opt(def)
	}

	return def
}
