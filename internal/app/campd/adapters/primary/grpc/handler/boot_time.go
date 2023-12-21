package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"github.com/awlsring/camp/pkg/metrics"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) BootTime(ctx context.Context, in *emptypb.Empty) (*campd.BootTimeOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("BootTime called")

	boot, err := h.hostSvc.BootTime(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get boot time")
		h.metrics.PutMetric("boot_time_fault", metrics.Counter, 1)
		return nil, errors.GrpcError(err)
	}
	bootTimestamp := model.NewTimestamp[uint64](boot)

	log.Debug().Msg("returning boot time")
	return &campd.BootTimeOutput{BootTime: bootTimestamp}, nil
}
