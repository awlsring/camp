package power_control

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *powerControlService) Reboot(ctx context.Context, identifier machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("Initiating reboot for machine %s", identifier)

	err := s.mRepo.UpdateStatus(ctx, identifier, machine.MachineStatusPending)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to set status of machine %s to pending", identifier)
		return err
	}

	log.Debug().Msgf("Sending reboot signal to machine %s", identifier)
	err = s.agent.RebootMachine(ctx, identifier.String(), endpoint.String(), token.String())
	if err != nil {
		log.Error().Err(err).Msgf("Failed to send reboot signal to machine %s", identifier)
		return err
	}

	log.Debug().Msgf("Setting status of machine %s to rebooting", identifier)
	err = s.mRepo.UpdateStatus(ctx, identifier, machine.MachineStatusRebooting)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to set status of machine %s to rebooting", identifier)
		return err
	}

	log.Debug().Msgf("Successfully sent reboot signal to machine %s", identifier)
	return nil
}
