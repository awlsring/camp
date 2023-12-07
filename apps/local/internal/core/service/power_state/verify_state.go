package power_state

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) VerifyState(ctx context.Context, id machine.Identifier, reported machine.MachineStatus, endpoint machine.MachineEndpoint, token machine.AgentKey) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("verifying state of %s", reported.String())

	log.Debug().Msg("checking connectivity for machine")
	actual, err := s.checkMachineConnectivity(ctx, id, endpoint, token)
	if err != nil {
		log.Error().Err(err).Msg("failed to check connectivity")
		return err
	}

	log.Debug().Msgf("actual state is %s, reported as %s", actual.String(), reported.String())

	if actual != reported {
		log.Debug().Msg("reported state does not match actual state")
		err := s.reportChangeAndUpdateState(ctx, id, reported, actual, false)
		if err != nil {
			return err
		}
	}

	return nil
}
