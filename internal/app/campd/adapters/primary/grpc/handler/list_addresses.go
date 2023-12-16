package handler

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/model"
	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) ListAddresses(ctx context.Context, in *emptypb.Empty) (*campd.ListAddressesOutput, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("listing addresses")

	addrs, err := h.netSvc.ListIpAddresses(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to list addresses")
		return nil, err
	}

	a := []*campd.IpAddressSummary{}
	for _, addr := range addrs {
		a = append(a, model.IpAddressFromDomain(addr))
	}

	log.Debug().Msg("returning addresses")
	return &campd.ListAddressesOutput{
		Addresses: a,
	}, nil
}
