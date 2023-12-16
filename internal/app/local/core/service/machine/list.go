package machine

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *machineService) ListMachines(ctx context.Context) ([]*machine.Machine, error) {
	log := logger.FromContext(ctx)

	log.Debug().Msg("Listing machines")
	machines, err := s.repo.List(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list machines")
		return nil, errors.New(errors.ErrInternalFailure, err)
	}

	log.Debug().Msg("Machines listed")
	return machines, nil
}
