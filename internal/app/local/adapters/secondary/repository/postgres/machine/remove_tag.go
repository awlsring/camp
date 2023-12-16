package machine_repository

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/tag"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (r *Repository) RemoveTag(ctx context.Context, id machine.Identifier, tag tag.TagKey) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Removing tag %s from machine %s", tag, id.String())

	err := r.tagRepo.DeleteTagFromResource(ctx, tag, id.String())
	if err != nil {
		log.Error().Err(err).Msgf("Failed to remove tag %s from machine %s", tag, id.String())
		return err
	}

	log.Debug().Msg("Successfully removed tag from machine")
	return nil
}
