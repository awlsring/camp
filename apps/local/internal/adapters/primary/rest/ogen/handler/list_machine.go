package handler

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) ListMachines(ctx context.Context) (camplocal.ListMachinesRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke ListMachines")
	m, err := h.mSvc.ListMachines(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list machines")
		return nil, err
	}
	log.Debug().Msgf("Found %d machines", len(m))

	var summaries []camplocal.MachineSummary
	for _, machine := range m {
		log.Debug().Msgf("Converting machine: %s", machine.Identifier.String())
		summaries = append(summaries, modelToSummary(machine))
	}

	return &camplocal.ListMachinesResponseContent{
		Summaries: summaries,
	}, nil
}
