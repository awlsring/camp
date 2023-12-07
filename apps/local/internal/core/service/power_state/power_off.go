package power_state

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) PowerOff(ctx context.Context, id machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Initiating power off for machine")

	log.Debug().Msg("Sending power off signal to machine %s")
	err := s.campd.PowerOffMachine(ctx, id.String(), endpoint.String(), token.String())
	if err != nil {
		log.Error().Err(err).Msg("Failed to send power off signal to machine")
		return err
	}

	log.Debug().Msg("Setting status of machine to pending")
	err = s.reportChangeAndUpdateState(ctx, id, machine.MachineStatusPending, machine.MachineStatusStopping, true)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set status of machine to pending")
		return err
	}

	log.Debug().Msg("Successfully sent power off to machine")
	return nil
}
