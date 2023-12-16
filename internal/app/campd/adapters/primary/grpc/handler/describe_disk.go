package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
)

func (h *Handler) DescribeDisk(ctx context.Context, in *campd.DescribeDiskInput) (*campd.DescribeDiskOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("DescribeDisk called")

	log.Debug().Msg("getting disk")
	disk, err := h.strSvc.DescribeDisk(ctx, in.Name)
	if err != nil {
		log.Error().Err(err).Msg("failed to get disk")
		return nil, err
	}

	log.Debug().Msg("returning Disk Summary")
	return &campd.DescribeDiskOutput{Disk: model.DiskFromDomain(disk)}, nil
}
