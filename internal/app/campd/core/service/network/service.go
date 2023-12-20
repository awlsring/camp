package network

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/network"
)

type ServiceOpt func(*Service)

func WithIgnoredPrefix(ignored []string) ServiceOpt {
	return func(i *Service) {
		i.ignoredPrefix = ignored
	}
}

type Service struct {
	loadVirtual   bool
	ignoredPrefix []string
	nics          map[string]*network.Nic
	addresses     []*network.IpAddress
}

func InitService(ctx context.Context, opts ...ServiceOpt) (service.Network, error) {
	s := &Service{
		loadVirtual:   false,
		ignoredPrefix: []string{"br", "veth", "docker", "cni", "flannel"},
	}

	for _, opt := range opts {
		opt(s)
	}

	if err := s.load(ctx); err != nil {
		return nil, err
	}

	return s, nil
}
