package storage

import (
	"context"
	"strings"
	"time"

	"github.com/awlsring/camp/internal/pkg/domain/storage"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/jaypipes/ghw"
)

func inIgnoreList(ignored []string, device string) bool {
	for _, ignore := range ignored {
		if strings.Contains(device, ignore) {
			return true
		}
	}

	return false
}

func (s *Service) refreshIfNeeded(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Checking if refresh is needed")
	now := time.Now().UTC()
	if now.Sub(s.lastCheck) > s.refreshInterval {
		log.Debug().Msg("Refresh needed, refreshing disks")
		err := s.loadDisks(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) loadDisks(ctx context.Context) error {
	devs, err := ghw.Block()
	if err != nil {
		return err
	}

	disks := map[string][]*storage.Disk{}
	for _, disk := range devs.Disks {
		if inIgnoreList(s.ignoredDevices, disk.Name) {
			continue
		}

		d := storage.NewDisk(
			disk.Name,
			disk.SizeBytes,
			storage.DriveTypeFromString(disk.DriveType.String()),
			storage.StorageControllerFromString(disk.StorageController.String()),
			disk.IsRemovable,
			disk.Vendor,
			disk.Model,
			disk.SerialNumber,
			disk.WWN,
		)
		disks[d.Name] = append(disks[d.Name], d)
	}

	return nil
}
