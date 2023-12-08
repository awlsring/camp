package power_state

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) ReconcileUnknownState(ctx context.Context, id machine.Identifier) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("machine is in an unknown state")

	log.Debug().Msg("getting machine entry")
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get machine entry")
		return err
	}

	log.Debug().Msgf("checking connectivity")
	state, err := s.checkMachineConnectivity(ctx, id, m.AgentEndpoint, m.AgentApiKey)
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
