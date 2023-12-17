package reporter

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) Heartbeat(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("heartbeating system")

	err := s.reporting.Heartbeat(ctx, s.id)
	if err != nil {
		log.Error().Err(err).Msg("failed to heartbeat system")
		return err
	}

	log.Info().Msg("heartbeating registered")
	return nil
}
