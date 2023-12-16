package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DescribeMemory(ctx context.Context, req *emptypb.Empty) (*campd.DescribeMemoryOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Describing memory")
	mem, err := h.memSvc.Description(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to describe memory")
		return nil, err
	}
	return &campd.DescribeMemoryOutput{
		Memory: model.MemoryFromDomain(mem),
	}, nil
}
