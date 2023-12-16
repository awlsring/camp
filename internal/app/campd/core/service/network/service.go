package network

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/network"
)

type ServiceOpt func(*Service)

type Service struct {
	loadVirtual bool
	nics        map[string]*network.Nic
}

func InitService(ctx context.Context, opts ...ServiceOpt) (service.Network, error) {
	s := &Service{
		loadVirtual: false,
	}

	for _, opt := range opts {
		opt(s)
	}

	if err := s.loadNics(ctx); err != nil {
		return nil, err
	}

	return s, nil
}
