package handler

import (
	"context"
	"strings"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/pkg/gen/local"
)

func (h *Handler) DescribeMachine(ctx context.Context, req camplocal.DescribeMachineParams) (camplocal.DescribeMachineRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Identifier: %s", req.Identifier)

	id, err := machine.IdentifierFromString(req.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.Identifier)
		return nil, err
	}

	m, err := h.mSvc.DescribeMachine(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Debug().Msgf("Machine with identifier %s not found", req.Identifier)
			return nil, err // TODO: Consolidate this error handling
		}
		log.Error().Err(err).Msgf("Failed to describe machine with identifier %s", req.Identifier)
		return nil, err

	}

	return &camplocal.DescribeMachineResponseContent{
		Summary: modelToSummary(m),
	}, nil
}
