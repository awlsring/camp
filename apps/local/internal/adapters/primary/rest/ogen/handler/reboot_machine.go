package handler

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) RebootMachine(ctx context.Context, req camplocal.RebootMachineParams) (camplocal.RebootMachineRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Invoking reboot machine for %s", req.Identifier)
	return &camplocal.RebootMachineResponseContent{
		Status: camplocal.MachineStatusPending,
	}, nil
}
