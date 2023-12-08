package power_state_job

import (
	"context"
	"encoding/json"

	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (j *RequestStateChangeJob) requestStateChangeMessageToDomain(ctx context.Context, msg []byte) (*power.RequestStateChangeMessage, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("unmarshalling message")

	var message power.RequestStateChangeMessageJson
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Error().Err(err).Msgf("failed to unmarshal message")
		return nil, err
	}

	log.Debug().Msgf("converting message to domain")
	domainMessage, err := message.ToDomain()
	if err != nil {
		log.Error().Err(err).Msgf("failed to convert message to domain")
		return nil, err
	}

	return domainMessage, nil
}

func (j *ValidateStateChangeJob) validateStateChangeMessageToDomain(ctx context.Context, msg []byte) (*power.StateValidationMessage, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("unmarshalling message")

	var message power.StateValidationMessageJson
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Error().Err(err).Msgf("failed to unmarshal message")
		return nil, err
	}

	log.Debug().Msgf("converting message to domain")
	domainMessage, err := message.ToDomain()
	if err != nil {
		log.Error().Err(err).Msgf("failed to convert message to domain")
		return nil, err
	}

	return domainMessage, nil
}
