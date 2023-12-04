package exchange

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ExchangeOpt func(*Definition)

type Definition struct {
	Name       string
	Type       ExchangeType
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

func (e Definition) CreateExchange(ctx context.Context, channel *amqp.Channel) error {
	log := logger.FromContext(ctx)

	log.Debug().Msgf("declaring exchange %s", e.Name)
	err := channel.ExchangeDeclare(e.Name, e.Type.String(), e.Durable, e.AutoDelete, e.Internal, e.NoWait, e.Args)
	if err != nil {
		log.Error().Err(err).Msgf("failed to declare exchange %s", e.Name)
		return err
	}

	log.Debug().Msgf("exchange %s declared", e.Name)
	return nil
}

func WithDurable(durable bool) ExchangeOpt {
	return func(e *Definition) {
		e.Durable = durable
	}
}

func WithAutoDelete(autoDelete bool) ExchangeOpt {
	return func(e *Definition) {
		e.AutoDelete = autoDelete
	}
}

func WithInternal(internal bool) ExchangeOpt {
	return func(e *Definition) {
		e.Internal = internal
	}
}

func WithNoWait(noWait bool) ExchangeOpt {
	return func(e *Definition) {
		e.NoWait = noWait
	}
}

func WithArgs(args amqp.Table) ExchangeOpt {
	return func(e *Definition) {
		e.Args = args
	}
}

func NewDefinition(name string, exchangeType ExchangeType, opts ...ExchangeOpt) *Definition {
	ex := &Definition{
		Name:       name,
		Type:       exchangeType,
		Durable:    false,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}

	for _, opt := range opts {
		opt(ex)
	}

	return ex
}
