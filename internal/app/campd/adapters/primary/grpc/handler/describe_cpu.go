package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DescribeCpu(ctx context.Context, in *emptypb.Empty) (*campd.DescribeCpuOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("DescribeCpu called")

	log.Debug().Msg("getting cpu")
	cpu, err := h.cpuSvc.Description(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get cpu")
		return nil, errors.GrpcError(err)
	}

	log.Debug().Msg("returning Cpu Summary")
	return &campd.DescribeCpuOutput{Cpu: model.CpuFromDomain(cpu)}, nil
}
