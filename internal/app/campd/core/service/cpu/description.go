package cpu

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (c *Service) Description(ctx context.Context) (*cpu.CPU, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Returning CPU description")
	return c.cpu, nil
}
