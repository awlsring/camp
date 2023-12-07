package power_state

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) ReconcileUnknownState(ctx context.Context, id machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("machine is in an unknown state")

	log.Debug().Msgf("checking connectivity")
	state, err := s.checkMachineConnectivity(ctx, id, endpoint, token)
	if err != nil {
		log.Error().Err(err).Msgf("error checking connectivity")
		return err
	}

	err = s.reportChangeAndUpdateState(ctx, id, state, state, false)
	if err != nil {
		return err
	}

	log.Debug().Msgf("machine is now in state %s", state.String())
	return nil
}
