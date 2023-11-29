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

func (h *Handler) describeMachineErrorHandler(ctx context.Context, err error) (camplocal.DescribeMachineRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Hadnling err: %s", err.Error())

	log.Debug().Msgf("Checking if err is camperror.Error")
	var campErr *camperror.Error
	if errors.As(err, &campErr) {
		log.Debug().Msgf("err is camperror.Error")
		e := campErr.CampError()
		log.Debug().Msgf("err is camperror.Error, is %s", e.Error())
		switch e {
		case camperror.ErrResourceNotFound:
			log.Debug().Msg("Returning resource not found")
			return &camplocal.ResourceNotFoundExceptionResponseContent{
				Message: err.Error(),
			}, nil
		case camperror.ErrValidation:
			return &camplocal.ValidationExceptionResponseContent{
				Message: err.Error(),
			}, nil
		}
	}

	log.Debug().Msg("Returing internal server error")
	return nil, err
}

func (h *Handler) DescribeMachine(ctx context.Context, req camplocal.DescribeMachineParams) (camplocal.DescribeMachineRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Identifier: %s", req.Identifier)

	id, err := machine.IdentifierFromString(req.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.Identifier)
		return h.describeMachineErrorHandler(ctx, err)
	}

	m, err := h.mSvc.DescribeMachine(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Debug().Msgf("Machine with identifier %s not found", req.Identifier)
			return h.describeMachineErrorHandler(ctx, err) // TODO: Consolidate this error handling
		}
		log.Error().Err(err).Msgf("Failed to describe machine with identifier %s", req.Identifier)
		return h.describeMachineErrorHandler(ctx, err)

	}

	return &camplocal.DescribeMachineResponseContent{
		Summary: modelToSummary(m),
	}, nil
}
