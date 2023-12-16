package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) ListDisks(ctx context.Context, in *emptypb.Empty) (*campd.ListDisksOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("ListDisks called")

	log.Debug().Msg("getting disks")
	disks, err := h.strSvc.ListDisks(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get disks")
		return nil, err
	}

	diskSummaries := make([]*campd.DiskSummary, len(disks))
	for i, disk := range disks {
		diskSummaries[i] = model.DiskFromDomain(disk)
	}

	log.Debug().Msg("returning Disk Summary")
	return &campd.ListDisksOutput{Disks: diskSummaries}, nil
}
