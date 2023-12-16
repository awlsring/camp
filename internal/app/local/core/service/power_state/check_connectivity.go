package power_state

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) checkMachineConnectivity(ctx context.Context, id machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) (power.StatusCode, error) {
	log := logger.FromContext(ctx)

	log.Debug().Msgf("checking connectivity")
	ok, err := s.campd.CheckMachineConnectivity(ctx, id.String(), endpoint.String(), token.String())
	if err != nil {
		log.Error().Err(err).Msgf("error checking connectivity")
		return power.StatusCodeStopping, err
	}

	log.Debug().Msgf("connectivity check returned %t", ok)
	if ok {
		return power.StatusCodeRunning, nil
	} else {
		return power.StatusCodeStopped, nil
	}
}
