package board

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/motherboard"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) DescribeMotherboard(ctx context.Context) (*motherboard.Motherboard, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Returning motherboard description")

	return s.mobo, nil
}
