package machine

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *machineService) AcknowledgeHeartbeat(ctx context.Context, id machine.Identifier) error {
	log := logger.FromContext(ctx)

	log.Debug().Msgf("Acknowledge heartbeat from machine %s", id)
	err := s.repo.UpdateHeartbeat(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to acknowledge heartbeat from machine %s", id)
		return errors.New(errors.ErrInternalFailure, err)
	}

	log.Debug().Msgf("Heartbeat from machine %s acknowledged", id)
	return nil
}
