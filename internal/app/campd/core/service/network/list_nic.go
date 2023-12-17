package network

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/logger"
	"golang.org/x/exp/maps"
)

func (s *Service) ListNics(ctx context.Context) ([]*network.Nic, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Listing nics")

	if s.nics == nil || len(s.nics) == 0 {
		log.Debug().Msg("No nics found")
		return nil, nil
	}

	return maps.Values[map[string]*network.Nic](s.nics), nil
}
