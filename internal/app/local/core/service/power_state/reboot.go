package power_state

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) Reboot(ctx context.Context, id machine.Identifier) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Initiating reboot for machine")

	log.Debug().Msg("getting machine entry")
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get machine entry")
		return err
	}

	log.Debug().Msg("Sending reboot signal to machine")
	err = s.campd.RebootMachine(ctx, m.Identifier.String(), m.AgentEndpoint.String(), m.AgentApiKey.String())
	if err != nil {
		log.Error().Err(err).Msg("Failed to send reboot signal to machine")
		return err
	}

	log.Debug().Msg("Setting status of machine to pending")
	err = s.reportChangeAndUpdateState(ctx, id, power.StatusCodePending, power.StatusCodeRebooting, true)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set status of machine to pending")
		return err
	}

	log.Debug().Msg("Successfully sent reboot signal to machine")
	return nil
}
