package power_state

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) reportChangeAndUpdateState(ctx context.Context, id machine.Identifier, reported, actual machine.MachineStatus, planned bool) error {
	log := logger.FromContext(ctx)

	log.Debug().Msg("reporting state change")
	err := s.publishStateChange(ctx, id, reported, actual, planned)
	if err != nil {
		log.Error().Err(err).Msg("failed to publish state change")
	}
	log.Debug().Msg("adjusting state in database")
	err = s.repo.UpdateStatus(ctx, id, actual)
	if err != nil {
		log.Error().Err(err).Msg("failed to update state in database")
		return err
	}

	return nil
}
