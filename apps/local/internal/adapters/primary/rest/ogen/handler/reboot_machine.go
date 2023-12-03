package handler

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) RebootMachine(ctx context.Context, req camplocal.RebootMachineParams) (camplocal.RebootMachineRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Invoking reboot machine for %s", req.Identifier)

	log.Debug().Msgf("Parsing machine identifier %s", req.Identifier)
	id, err := machine.IdentifierFromString(req.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse machine identifier %s", req.Identifier)
		return nil, err
	}

	log.Debug().Msgf("Sending reboot request for machine %s", id)
	err = h.mSvc.RequestPowerChange(ctx, id, power.ChangeTypeReboot)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Successfully sent reboot request for machine %s, reporting as machine status pending", id)
	return &camplocal.RebootMachineResponseContent{
		Status: camplocal.MachineStatusPending,
	}, nil
}
