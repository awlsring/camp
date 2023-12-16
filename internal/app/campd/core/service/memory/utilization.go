package memory

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/shirou/gopsutil/mem"
)

func (s *Service) Utilization(ctx context.Context) (*memory.Utilization, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Getting current memory utilization")

	v, err := mem.VirtualMemory()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current memory utilization")
		return nil, err
	}

	log.Debug().Msg("Successfully got current memory utilization")
	return &memory.Utilization{
		Total: v.Total,
		Used:  v.Used,
		Free:  v.Free,
	}, nil
}
