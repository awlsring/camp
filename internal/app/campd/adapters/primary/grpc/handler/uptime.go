package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) Uptime(ctx context.Context, in *emptypb.Empty) (*campd.UptimeOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Uptime called")

	log.Debug().Msg("getting uptime")
	up, err := h.hostSvc.Uptime(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get uptime")
		return nil, errors.GrpcError(err)
	}
	upTimestamp := model.NewTimestamp[uint64](up)

	log.Debug().Msg("returning boot time")
	return &campd.UptimeOutput{Uptime: upTimestamp}, nil
}
