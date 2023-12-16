package power_state

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) PowerOff(ctx context.Context, id machine.Identifier) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Initiating power off for machine")

	log.Debug().Msg("getting machine entry")
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get machine entry")
		return err
	}

	log.Debug().Msg("Sending power off signal to machine %s")
	err = s.campd.PowerOffMachine(ctx, m.Identifier.String(), m.AgentEndpoint.String(), m.AgentApiKey.String())
	if err != nil {
		log.Error().Err(err).Msg("Failed to send power off signal to machine")
		return err
	}

	log.Debug().Msg("Setting status of machine to pending")
	err = s.reportChangeAndUpdateState(ctx, id, power.StatusCodePending, power.StatusCodeStopping, true)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set status of machine to pending")
		return err
	}

	log.Debug().Msg("Successfully sent power off to machine")
	return nil
}
