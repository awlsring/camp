package machine

import (
	"context"
	"errors"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
	"github.com/awlsring/camp/apps/local/internal/ports/repository"
	camperrors "github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *machineService) AddTagsToMachine(ctx context.Context, id machine.Identifier, tags []*tag.Tag) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Adding tags to machine %s", id.String())

	err := s.repo.AddTags(ctx, id, tags)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateResourceTag) {
			log.Warn().Err(err).Msgf("Failed to add tags to machine %s", id)
			return camperrors.New(camperrors.ErrDuplicate, errors.New("request contains tag key that already exists on resource"))
		}
		log.Error().Err(err).Msgf("Failed to add all tags to machine %s", id)
		return camperrors.New(camperrors.ErrInternalFailure, err)
	}

	log.Debug().Msgf("%d tags added to machine %s", len(tags), id)
	return nil
}
