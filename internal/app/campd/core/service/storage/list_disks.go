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
		log.Error().Err(err).Msg("error refreshing disks")
		return nil, err
	}

	if s.disks == nil || len(s.disks) == 0 {
		log.Debug().Msg("No disks found")
		return nil, nil
	}

	disk := maps.Values[map[string]*storage.Disk](s.disks)
	return disk, nil
}
