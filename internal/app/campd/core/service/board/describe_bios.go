package board

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/motherboard"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) DescribeBios(ctx context.Context) (*motherboard.Bios, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Returning BIOS description")

	return s.bios, nil
}
