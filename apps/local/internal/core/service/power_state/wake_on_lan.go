package power_state

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) WakeOnLan(ctx context.Context, id machine.Identifier) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("Initiating wake on lan for machine %s", id)

	log.Debug().Msg("getting machine entry")
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get machine entry")
		return err
	}

	log.Debug().Msg("Sending magic packet")
	err = s.wol.SendSignal(ctx, m.PowerCapabilities.WakeOnLan.MacAddress.String())
	if err != nil {
		log.Error().Err(err).Msg("Failed to send magic packet")
		return errors.New(errors.ErrInternalFailure, err)
	}
	log.Debug().Msg("Successfully sent magic packet")

	log.Debug().Msg("Setting status of machine to pending")
	err = s.reportChangeAndUpdateState(ctx, id, machine.MachineStatusPending, machine.MachineStatusStarting, true)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to set status of machine %s to pending", id)
		return errors.New(errors.ErrInternalFailure, err)
	}

	return nil
}
