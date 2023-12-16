package storage

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/storage"
	"github.com/awlsring/camp/internal/pkg/logger"
	"golang.org/x/exp/maps"
)

func (s *Service) ListDisks(ctx context.Context) ([]*storage.Disk, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Listing disks")
	err := s.refreshIfNeeded(ctx)
	if err != nil {
		return nil, err
	}

	return maps.Values[map[string]*storage.Disk](s.disks), nil
}
