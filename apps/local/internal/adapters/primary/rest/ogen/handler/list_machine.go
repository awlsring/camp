package handler

import (
	"context"
	"errors"

	camperror "github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) listMachinesErrorHandler(err error) (camplocal.ListMachinesRes, error) {
	var campErr *camperror.Error
	if errors.As(err, &campErr) {
		e := campErr.CampError()
		switch e {
		case camperror.ErrValidation:
			return &camplocal.ValidationExceptionResponseContent{
				Message: err.Error(),
			}, nil
		}
	}
	return nil, err
}

func (h *Handler) ListMachines(ctx context.Context) (camplocal.ListMachinesRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke ListMachines")
	m, err := h.mSvc.ListMachines(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list machines")
		return h.listMachinesErrorHandler(err)
	}
	log.Debug().Msgf("Found %d machines", len(m))

	var summaries []camplocal.MachineSummary
	for _, machine := range m {
		log.Debug().Msgf("Converting machine: %+v", machine)
		summaries = append(summaries, modelToSummary(machine))
	}

	return &camplocal.ListMachinesResponseContent{
		Summaries: summaries,
	}, nil
}
