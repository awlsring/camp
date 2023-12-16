package host

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/shirou/gopsutil/v3/host"
)

func (s *Service) BootTime(ctx context.Context) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Returning host boot time")

	t, err := host.BootTimeWithContext(ctx)
	if err != nil {
		return 0, err
	}
	return t, nil
}
