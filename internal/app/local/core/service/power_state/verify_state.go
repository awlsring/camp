package power_state

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) VerifyFinalState(ctx context.Context, id machine.Identifier, reported power.StatusCode) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("verifying state of %s", reported.String())

	log.Debug().Msg("getting machine entry")
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get machine entry")
		return err
	}

	err = s.verifyState(ctx, m, reported)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) verifyState(ctx context.Context, m *machine.Machine, reported power.StatusCode) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("verifying state of %s", reported.String())

	log.Debug().Msg("checking connectivity for machine")
	actual, err := s.checkMachineConnectivity(ctx, m.Identifier, m.AgentEndpoint, m.AgentApiKey)
	if err != nil {
		log.Error().Err(err).Msg("failed to check connectivity")
		return err
	}

	log.Debug().Msgf("actual state is %s, reported as %s", actual.String(), reported.String())

	if actual != reported {
		log.Debug().Msg("reported state does not match actual state")
		err := s.reportChangeAndUpdateState(ctx, m.Identifier, reported, actual, false)
		if err != nil {
			return err
		}
	}

	return nil
}
