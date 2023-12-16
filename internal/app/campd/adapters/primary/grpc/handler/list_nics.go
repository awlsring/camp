package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) ListNetworkInterfaces(ctx context.Context, in *emptypb.Empty) (*campd.ListNetworkInterfacesOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("listing network interfaces")

	nics, err := h.netSvc.ListNics(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to list network interfaces")
		return nil, err
	}

	n := []*campd.NetworkInterfaceSummary{}
	for _, nic := range nics {
		n = append(n, model.NicFromDomain(nic))
	}

	log.Debug().Msg("returning network interfaces")
	return &campd.ListNetworkInterfacesOutput{
		Nics: n,
	}, nil
}
