package power_state_topic

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (t *PowerStateExchange) SendRequestStateChangeMessage(ctx context.Context, msg *power.RequestStateChangeMessage) error {
	log := logger.FromContext(ctx)

	log.Debug().Msg("converting power change request message to json")
	b, err := msg.ToJson()
	if err != nil {
		return err
	}

	log.Debug().Msg("publishing power change request message")
	err = t.publish(ctx, b, RequestStateChangeRequestKey)
	if err != nil {
		return err
	}

	log.Debug().Msg("power change request message published")
	return nil
}
