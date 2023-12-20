package network

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/network"
)

type ServiceOpt func(*Service)

func WithIgnoredNicPrefix(ignored []string) ServiceOpt {
	return func(i *Service) {
		i.ignoredNicPrefix = ignored
	}
}

func WithIgnoredAddrPrefix(ignored []string) ServiceOpt {
	return func(i *Service) {
		i.ignoredAddrPrefix = ignored
	}
}

type Service struct {
	loadVirtual       bool
	ignoredNicPrefix  []string
	ignoredAddrPrefix []string
	nics              map[string]*network.Nic
	addresses         []*network.IpAddress
}

func InitService(ctx context.Context, opts ...ServiceOpt) (service.Network, error) {
	s := &Service{
		loadVirtual:       false,
		ignoredNicPrefix:  []string{"br", "veth", "docker", "cni", "flannel"},
		ignoredAddrPrefix: []string{"fe80::", "169.254."},
	}

	for _, opt := range opts {
		opt(s)
	}

	if err := s.load(ctx); err != nil {
		return nil, err
	}

	return s, nil
}
