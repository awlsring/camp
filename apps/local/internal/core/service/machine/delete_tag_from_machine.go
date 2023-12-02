package machine

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *machineService) RemoveTagFromMachine(ctx context.Context, id machine.Identifier, tag tag.TagKey) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Removing tag %s from machine %s", tag, id)

	err := s.repo.RemoveTag(ctx, id, tag)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to remove tag %s from machine %s", tag, id)
		return errors.New(errors.ErrInternalFailure, err)
	}

	return nil
}
