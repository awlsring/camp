package handler

import (
	"context"
	"errors"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	camperror "github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) reportStatusChangeErrorHandler(err error) (camplocal.ReportStatusChangeRes, error) {
	var campErr camperror.Error
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

func (h *Handler) ReportStatusChange(ctx context.Context, req *camplocal.ReportStatusChangeRequestContent) (camplocal.ReportStatusChangeRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke ReportStatusChange")

	iid, err := machine.InternalIdentifierFromString(req.InternalIdentifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.InternalIdentifier)
		return h.reportStatusChangeErrorHandler(err)
	}

	status, err := machine.MachineStatusFromString(string(req.Status))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse status %s", req.Status)
		return h.reportStatusChangeErrorHandler(err)
	}

	err = h.mSvc.UpdateStatus(ctx, iid, status)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update machine status")
		return h.reportStatusChangeErrorHandler(err)
	}
	return &camplocal.ReportStatusChangeResponseContent{
		Success: true,
	}, nil
}
