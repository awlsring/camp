package handler

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func (h *Handler) Heartbeat(ctx context.Context, req *camplocal.HeartbeatRequestContent) (camplocal.HeartbeatRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke Heartbeat")

	id, err := machine.IdentifierFromString(req.InternalIdentifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse identifier %s", req.InternalIdentifier)
		return nil, err
	}

	err = h.mSvc.AcknowledgeHeartbeat(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to acknowledge heartbeat")
		return nil, err
	}

	return &camplocal.HeartbeatResponseContent{
		Success: true,
	}, nil
}
