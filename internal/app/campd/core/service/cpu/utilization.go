package cpu

import (
	"context"

	pscpu "github.com/shirou/gopsutil/v3/cpu"

	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (c *Service) Utilization(ctx context.Context) ([]*cpu.Utilization, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Getting CPU utilization")

	percent, err := pscpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get CPU utilization")
		return nil, err
	}

	utilizations := make([]*cpu.Utilization, len(percent))
	for i, p := range percent {

		utilizations[i] = &cpu.Utilization{
			Core:  i,
			Usage: p,
		}
	}
	return utilizations, nil
}
