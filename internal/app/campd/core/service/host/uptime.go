package host

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/shirou/gopsutil/v3/host"
)

func (s *Service) Uptime(ctx context.Context) (uint64, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Returning host uptime")

	t, err := host.UptimeWithContext(ctx)
	if err != nil {
		return 0, err
	}
	return t, nil
}
