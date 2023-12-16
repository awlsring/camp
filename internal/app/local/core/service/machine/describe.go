package machine

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/app/local/ports/repository"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *machineService) DescribeMachine(ctx context.Context, identifier machine.Identifier) (*machine.Machine, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Describing machine with identifier %s", identifier)

	m, err := s.repo.Get(ctx, identifier)
	if err != nil {
		if err == repository.ErrMachineNotFound {
			log.Debug().Err(err).Msgf("Machine with identifier %s not found", identifier)
			return nil, errors.New(errors.ErrResourceNotFound, err)
		}
		log.Error().Err(err).Msgf("Failed to describe machine with identifier %s", identifier)
		return nil, errors.New(errors.ErrInternalFailure, err)
	}

	log.Debug().Msgf("Machine with identifier %s found", identifier)
	return m, nil
}
