package network

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) ListIpAddresses(ctx context.Context) ([]*network.IpAddress, error) {
	log := logger.FromContext(ctx)
	log.Warn().Msg("Ip addressing not implemented, returning nil")
	return nil, nil
}
