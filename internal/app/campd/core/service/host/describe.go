package host

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) Describe(ctx context.Context) (*host.Host, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Returning host description")

	return s.host, nil
}
