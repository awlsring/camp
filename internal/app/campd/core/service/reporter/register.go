package reporter

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) Register(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("registering system")

	return nil
}
