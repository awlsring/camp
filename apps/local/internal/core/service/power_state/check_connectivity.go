package power_state

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) checkMachineConnectivity(ctx context.Context, id machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) (machine.MachineStatus, error) {
	log := logger.FromContext(ctx)

	log.Debug().Msgf("checking connectivity")
	ok, err := s.campd.CheckMachineConnectivity(ctx, id.String(), endpoint.String(), token.String())
	if err != nil {
		log.Error().Err(err).Msgf("error checking connectivity")
		return machine.MachineStatusUnknown, err
	}

	log.Debug().Msgf("connectivity check returned %t", ok)
	if ok {
		return machine.MachineStatusRunning, nil
	} else {
		return machine.MachineStatusStopped, nil
	}
}
