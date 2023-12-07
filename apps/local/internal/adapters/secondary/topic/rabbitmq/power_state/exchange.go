package power_state_topic

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/ports/topic"
	"github.com/awlsring/camp/internal/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName                 = "power_state"
	ExchangeTypeTopic            = "topic"
	RequestStateChangeRequestKey = "power_state.request_state_change"
	StateValidationRequestKey    = "power_state.state_validation"
)

var _ topic.PowerStateJob = &PowerStateExchange{}

type PowerStateExchange struct {
	channel *amqp.Channel
}

func New(channel *amqp.Channel) topic.PowerStateJob {
	return &PowerStateExchange{
		channel: channel,
	}
}

func (t *PowerStateExchange) Init(ctx context.Context, channel *amqp.Channel) error {
	log := logger.FromContext(ctx)

	log.Debug().Msg("initializing power change exchange")
	err := t.channel.ExchangeDeclare(ExchangeName, ExchangeTypeTopic, false, false, false, true, nil)
	if err != nil {
		return err
	}

	return nil
}

func (t *PowerStateExchange) publish(ctx context.Context, msg []byte, key string) error {
	log := logger.FromContext(ctx)

	log.Debug().Msg("forming ampq message")
	ampMsg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         msg,
	}

	log.Debug().Msgf("publishing ampq message to excahnge %s", ExchangeName)
	err := t.channel.PublishWithContext(ctx, ExchangeName, key, true, false, ampMsg)
	if err != nil {
		return err
	}

	return nil
}
