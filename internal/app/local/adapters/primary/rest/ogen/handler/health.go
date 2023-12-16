package handler

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
	camplocal "github.com/awlsring/camp/pkg/gen/local"
)

func (h *Handler) Health(ctx context.Context) (camplocal.HealthRes, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke Health")
	return &camplocal.HealthResponseContent{
		Success: true,
	}, nil
}
