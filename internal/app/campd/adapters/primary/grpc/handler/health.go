package handler

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) Health(ctx context.Context, in *emptypb.Empty) (*campd.HealthOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("health check")
	return &campd.HealthOutput{
		Success: true,
	}, nil
}
