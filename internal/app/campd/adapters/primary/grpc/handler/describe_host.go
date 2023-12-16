package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DescribeHost(ctx context.Context, req *emptypb.Empty) (*campd.DescribeHostOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Describing host")
	host, err := h.hostSvc.Describe(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to describe host")
		return nil, errors.GrpcError(err)
	}

	log.Debug().Msg("Returning host summary")
	return &campd.DescribeHostOutput{
		Host: model.HostFromDomain(host),
	}, nil
}
