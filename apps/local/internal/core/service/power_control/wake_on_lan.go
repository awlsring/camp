package power_control

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *powerControlService) WakeOnLan(ctx context.Context, id machine.Identifier, mac machine.MacAddress) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("Initiating wake on lan for machine %s", id)

	err := s.mRepo.UpdateStatus(ctx, id, machine.MachineStatusPending)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to set status of machine %s to pending", id)
		return errors.New(errors.ErrInternalFailure, err)
	}

	log.Debug().Msgf("Sending magic packet to %s", mac)
	err = s.wol.SendSignal(ctx, mac.String())
	if err != nil {
		log.Error().Err(err).Msgf("Failed to send magic packet to %s", mac)
		return errors.New(errors.ErrInternalFailure, err)
	}
	log.Debug().Msgf("Successfully sent magic packet to %s", mac)

	log.Debug().Msgf("Setting status of machine %s to starting", id)
	err = s.mRepo.UpdateStatus(ctx, id, machine.MachineStatusStarting)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to set status of machine %s to starting", id)
		return errors.New(errors.ErrInternalFailure, err)
	}

	return nil
}
