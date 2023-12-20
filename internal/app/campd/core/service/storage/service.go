package storage

import (
	"context"
	"time"

	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/storage"
)

const DefaultRefreshInterval = 6 * time.Hour

type ServiceOpt func(*Service)

func WithIgnoredDevices(ignored []string) ServiceOpt {
	return func(i *Service) {
		i.ignoredDevices = ignored
	}
}

func WithRefreshInterval(interval time.Duration) ServiceOpt {
	return func(i *Service) {
		i.refreshInterval = interval
	}
}

type Service struct {
	refreshInterval time.Duration
	lastCheck       time.Time
	ignoredDevices  []string
	disks           map[string]*storage.Disk
}

func InitService(ctx context.Context, opts ...ServiceOpt) (service.Storage, error) {
	s := &Service{
		refreshInterval: DefaultRefreshInterval,
		ignoredDevices: []string{
			"dm", "zram",
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	if err := s.loadDisks(ctx); err != nil {
		return nil, err
	}
	s.lastCheck = time.Now().UTC()

	return s, nil
}
