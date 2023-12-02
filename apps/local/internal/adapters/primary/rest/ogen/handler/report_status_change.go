package handler

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) ReportStatusChange(ctx context.Context, req *camplocal.ReportStatusChangeRequestContent) (camplocal.ReportStatusChangeRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke ReportStatusChange")

	id, err := machine.IdentifierFromString(req.InternalIdentifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.InternalIdentifier)
		return nil, err
	}

	status, err := machine.MachineStatusFromString(string(req.Status))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse status %s", req.Status)
		return nil, err
	}

	err = h.mSvc.UpdateStatus(ctx, id, status)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update machine status")
		return nil, err
	}
	return &camplocal.ReportStatusChangeResponseContent{
		Success: true,
	}, nil
}
