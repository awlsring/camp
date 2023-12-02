package handler

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) PowerOffMachine(ctx context.Context, req camplocal.PowerOffMachineParams) (camplocal.PowerOffMachineRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Invoking power off machine for %s", req.Identifier)
	return &camplocal.PowerOffMachineResponseContent{
		Success: true,
	}, nil
}
