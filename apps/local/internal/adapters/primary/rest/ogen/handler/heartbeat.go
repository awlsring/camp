package handler

import (
	"context"
	"errors"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	camperror "github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) hearbeatErrorHandler(err error) (camplocal.HeartbeatRes, error) {
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

func (h *Handler) Heartbeat(ctx context.Context, req *camplocal.HeartbeatRequestContent) (camplocal.HeartbeatRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke Heartbeat")

	id, err := machine.IdentifierFromString(req.InternalIdentifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.InternalIdentifier)
		return h.hearbeatErrorHandler(err)
	}

	err = h.mSvc.AcknowledgeHeartbeat(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to acknowledge heartbeat")
		return h.hearbeatErrorHandler(err)
	}

	return &camplocal.HeartbeatResponseContent{
		Success: true,
	}, nil
}
