package power_change_topic

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
	"github.com/awlsring/camp/apps/local/internal/ports/topic"
	"github.com/awlsring/camp/internal/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName          = "power_change"
	ExchangeTypeTopic     = "topic"
	PowerChangeRequestKey = "power_change.request"
)

var _ topic.PowerChange = &PowerChangeExchange{}

type PowerChangeExchange struct {
	channel *amqp.Channel
}

func New() topic.PowerChange {
	return &PowerChangeExchange{}
}

func (t *PowerChangeExchange) Init(ctx context.Context, channel *amqp.Channel) error {
	log := logger.FromContext(ctx)

	log.Debug().Msg("initializing power change exchange")
	err := t.channel.ExchangeDeclare(ExchangeName, ExchangeTypeTopic, false, false, false, true, nil)
	if err != nil {
		return err
	}

	return nil
}

func (t *PowerChangeExchange) SendPowerChangeRequest(ctx context.Context, msg *power.PowerChangeRequestMessage) error {
	log := logger.FromContext(ctx)

	jsonMsg := PowerChangeRequestMessageJsonFromDomain(msg)
	b, err := jsonMsg.ToJson()
	if err != nil {
		return err
	}

	log.Debug().Msg("publishing power change request message")
	err = t.publish(ctx, b, PowerChangeRequestKey)
	if err != nil {
		return err
	}

	return nil
}

func (t *PowerChangeExchange) SendValidatePowerChangeRequest(ctx context.Context, msg *power.ValidatePowerChangeRequestMessage) error {
	log := logger.FromContext(ctx)

	log.Debug().Msg("validating power change request message")
	jsonMsg := ValidatePowerChangeRequestMessageJsonFromDomain(msg)
	b, err := jsonMsg.ToJson()
	if err != nil {
		return err
	}

	log.Debug().Msg("publishing power change validation request message")
	err = t.publish(ctx, b, PowerChangeRequestKey)
	if err != nil {
		return err
	}

	log.Debug().Msg("power change validation request message published")
	return nil
}

func (t *PowerChangeExchange) publish(ctx context.Context, msg []byte, key string) error {
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
