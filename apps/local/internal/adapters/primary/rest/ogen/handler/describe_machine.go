package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	camperror "github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) describeMachineErrorHandler(err error) (camplocal.DescribeMachineRes, error) {
	var campErr *camperror.Error
	if errors.As(err, &campErr) {
		e := campErr.CampError()
		switch e {
		case camperror.ErrResourceNotFound:
			return &camplocal.ResourceNotFoundExceptionResponseContent{
				Message: err.Error(),
			}, nil
		case camperror.ErrValidation:
			return &camplocal.ValidationExceptionResponseContent{
				Message: err.Error(),
			}, nil
		}
	}
	return nil, err
}

func (h *Handler) DescribeMachine(ctx context.Context, req camplocal.DescribeMachineParams) (camplocal.DescribeMachineRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Identifier: %s", req.Identifier)

	id, err := machine.IdentifierFromString(req.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.Identifier)
		return h.describeMachineErrorHandler(err)
	}

	m, err := h.mSvc.DescribeMachine(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Debug().Msgf("Machine with identifier %s not found", req.Identifier)
			return h.describeMachineErrorHandler(err) // TODO: Consolidate this error handling
		}
		log.Error().Err(err).Msgf("Failed to describe machine with identifier %s", req.Identifier)
		return h.describeMachineErrorHandler(err)

	}

	return &camplocal.DescribeMachineResponseContent{
		Summary: modelToSummary(m),
	}, nil
}
