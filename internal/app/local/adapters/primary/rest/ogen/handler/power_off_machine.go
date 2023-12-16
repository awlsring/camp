package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/app/local/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/pkg/gen/local"
)

func (h *Handler) PowerOffMachine(ctx context.Context, req camplocal.PowerOffMachineParams) (camplocal.PowerOffMachineRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Invoking power off machine for %s", req.Identifier)

	log.Debug().Msgf("Parsing machine identifier %s", req.Identifier)
	id, err := machine.IdentifierFromString(req.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse machine identifier %s", req.Identifier)
		return nil, err
	}

	log.Debug().Msgf("Sending power off request for machine %s", id)
	err = h.mSvc.RequestPowerChange(ctx, id, power.ChangeTypePowerOff)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Successfully sent power off request for machine %s, reporting as machine status pending", id)
	return &camplocal.PowerOffMachineResponseContent{
		Status: camplocal.StatusCodePending,
	}, nil
}
