package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DescribeMotherboard(ctx context.Context, req *emptypb.Empty) (*campd.DescribeMotherboardOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Describing motherboard")
	mb, err := h.moboSvc.DescribeMotherboard(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to describe motherboard")
		return nil, err
	}
	log.Debug().Msg("Returning motherboard summary")
	return &campd.DescribeMotherboardOutput{
		Motherboard: model.MotherboardFromDomain(mb),
	}, nil
}
