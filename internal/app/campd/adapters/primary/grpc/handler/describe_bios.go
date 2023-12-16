package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DescribeBIOS(ctx context.Context, in *emptypb.Empty) (*campd.DescribeBiosOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("DescribeBios called")

	log.Debug().Msg("getting bios")
	bios, err := h.moboSvc.DescribeBios(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get bios")
		return nil, errors.GrpcError(err)
	}

	log.Debug().Msg("returning Bios Summary")
	return &campd.DescribeBiosOutput{Bios: model.BiosFromDomain(bios)}, nil
}
