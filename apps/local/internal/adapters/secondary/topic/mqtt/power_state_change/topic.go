package power_state_change

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
	"github.com/awlsring/camp/apps/local/internal/ports/topic"
	"github.com/awlsring/camp/internal/pkg/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _ topic.PowerStateChange = &Topic{}

type Topic struct {
	topic  string
	client mqtt.Client
}

func New(client mqtt.Client, topic string) *Topic {
	return &Topic{
		topic:  topic,
		client: client,
	}
}

func (t *Topic) SendStateChangeMessage(ctx context.Context, msg *power.StateChangeMessage) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("sending power state change event")

	log.Debug().Msg("converting power state change event to json")
	b, err := msg.ToJson()
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal power state change event")
		return err
	}

	log.Debug().Msgf("sending power state change event: %s", b)
	token := t.client.Publish(t.topic, 0, false, b)
	if token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("failed to send power state change event")
		return token.Error()
	}
	log.Debug().Msg("sent power state change event")
	return nil
}
