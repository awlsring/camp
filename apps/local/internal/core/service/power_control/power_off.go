package power_control

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *powerControlService) PowerOff(ctx context.Context, identifier machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("Initiating power off for machine %s", identifier)

	err := s.mRepo.UpdateStatus(ctx, identifier, machine.MachineStatusPending)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to set status of machine %s to pending", identifier)
		return err
	}

	log.Debug().Msgf("Sending power off signal to machine %s", identifier)
	err = s.agent.PowerOffMachine(ctx, identifier.String(), endpoint.String(), token.String())
	if err != nil {
		log.Error().Err(err).Msgf("Failed to send power off signal to machine %s", identifier)
		return err
	}
	log.Debug().Msgf("Successfully sent power off signal to machine %s", identifier)

	log.Debug().Msgf("Setting status of machine %s to stopped", identifier)
	err = s.mRepo.UpdateStatus(ctx, identifier, machine.MachineStatusStopped)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to set status of machine %s to stopped", identifier)
		return err
	}

	log.Debug().Msgf("Successfully sent power off to machine %s", identifier)
	return nil
}
