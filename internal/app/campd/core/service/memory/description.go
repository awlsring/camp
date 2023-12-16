package memory

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) Description(ctx context.Context) (*memory.Memory, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Returning memory description")

	return &memory.Memory{
		Total: s.total,
	}, nil
}
