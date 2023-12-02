package machine

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *machineService) AddTagsToMachine(ctx context.Context, id machine.Identifier, tags []*tag.Tag) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Adding tags to machine %s", id.String())

	err := s.repo.AddTags(ctx, id, tags)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to add all tags to machine %s", id)
		return errors.New(errors.ErrInternalFailure, err)
	}

	log.Debug().Msgf("%d tags added to machine %s", len(tags), id)
	return nil
}
