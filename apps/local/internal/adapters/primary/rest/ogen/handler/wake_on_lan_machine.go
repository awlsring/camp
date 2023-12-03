package handler

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) SendMachineWakeOnLan(ctx context.Context, req camplocal.SendMachineWakeOnLanParams) (camplocal.SendMachineWakeOnLanRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Invoking wake on lan handler for machine %s", req.Identifier)

	log.Debug().Msgf("Parsing machine identifier %s", req.Identifier)
	id, err := machine.IdentifierFromString(req.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse machine identifier %s", req.Identifier)
		return nil, err
	}

	log.Debug().Msgf("Sending wake on lan request for machine %s", id)
	err = h.mSvc.RequestPowerChange(ctx, id, power.ChangeTypeWakeOnLan)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Successfully sent wake on lan request for machine %s, reporting as machine status pending", id)
	return &camplocal.SendMachineWakeOnLanResponseContent{
		Status: camplocal.MachineStatusPending,
	}, nil
}
