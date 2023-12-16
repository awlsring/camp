package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
)

func (h *Handler) DescribeNetworkInterface(ctx context.Context, in *campd.DescribeNetworkInterfaceInput) (*campd.DescribeNetworkInterfaceOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("DescribeNetworkInterface called")

	log.Debug().Msg("getting nic")
	nic, err := h.netSvc.DescribeNic(ctx, in.Name)
	if err != nil {
		log.Error().Err(err).Msg("failed to get nic")
		return nil, err
	}

	log.Debug().Msg("returning Nic Summary")
	return &campd.DescribeNetworkInterfaceOutput{Nic: model.NicFromDomain(nic)}, nil
}
