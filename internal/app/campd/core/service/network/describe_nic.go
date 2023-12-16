package network

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) DescribeNic(ctx context.Context, n string) (*network.Nic, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Describing nic")
	nic, ok := s.nics[n]
	if !ok {
		log.Debug().Msg("nic not found")
		return nil, service.ErrNicNotFound
	}
	log.Debug().Msg("returning nic")
	return nic, nil
}
