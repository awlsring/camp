package storage

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/storage"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) DescribeDisk(ctx context.Context, name string) (*storage.Disk, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Listing disks")
	err := s.refreshIfNeeded(ctx)
	if err != nil {
		return nil, err
	}

	disk, ok := s.disks[name]
	if !ok {
		return nil, service.ErrDiskNotFound
	}

	return disk, nil
}
