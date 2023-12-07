package power_state

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) publishStateChange(ctx context.Context, id machine.Identifier, reported, actual machine.MachineStatus, planned bool) error {
	log := logger.FromContext(ctx)

	log.Debug().Msgf("Creating state change message")
	msg := power.NewStateChangeMessage(id, reported, actual, planned)

	log.Debug().Msgf("Sending state change message")
	err := s.changeTopic.SendStateChangeMessage(ctx, msg)
	if err != nil {
		log.Error().Err(err).Msgf("failed to send state change message")
		return err
	}

	return nil
}
