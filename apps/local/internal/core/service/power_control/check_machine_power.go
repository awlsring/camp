package power_control

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *powerControlService) CheckMachinePower(ctx context.Context, identifier machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) (bool, error) {
	log := logger.FromContext(ctx)
	log.Info().Msgf("Initiating power off for machine %s", identifier)

	log.Debug().Msg("Sending a check to see if response is recieved from the machine")
	responds, err := s.agent.CheckMachineConnectivity(ctx, identifier.String(), endpoint.String(), token.String())
	if err != nil {
		log.Error().Err(err).Msgf("Failed to send check to machine %s", identifier)
		return false, err
	}

	return responds, nil
}
