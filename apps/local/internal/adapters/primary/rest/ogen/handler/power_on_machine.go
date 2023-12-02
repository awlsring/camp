package handler

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) PowerOnMachine(ctx context.Context, req camplocal.PowerOnMachineParams) (camplocal.PowerOnMachineRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Invoking power on machine for %s", req.Identifier)
	return &camplocal.PowerOnMachineResponseContent{
		Success: true,
	}, nil
}
