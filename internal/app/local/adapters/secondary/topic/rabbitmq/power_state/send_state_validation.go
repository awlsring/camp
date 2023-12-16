package power_state_topic

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (t *PowerStateExchange) SendStateValidationMessage(ctx context.Context, msg *power.StateValidationMessage) error {
	log := logger.FromContext(ctx)

	log.Debug().Msg("validating power change request message")
	b, err := msg.AsJson()
	if err != nil {
		return err
	}

	log.Debug().Msg("publishing power change validation request message")
	err = t.publish(ctx, b, StateValidationRequestKey)
	if err != nil {
		return err
	}

	log.Debug().Msg("power change validation request message published")
	return nil
}
