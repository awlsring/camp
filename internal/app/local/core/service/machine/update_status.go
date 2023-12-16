package machine

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *machineService) UpdateStatus(ctx context.Context, id machine.Identifier, status power.StatusCode) error {
	log := logger.FromContext(ctx)

	log.Debug().Msgf("Updating status of machine %s to %s", id, status)
	err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to update status of machine %s to %s", id, status)
		return errors.New(errors.ErrInternalFailure, err)
	}

	log.Debug().Msgf("Status of machine %s updated to %s", id, status)
	return nil
}
