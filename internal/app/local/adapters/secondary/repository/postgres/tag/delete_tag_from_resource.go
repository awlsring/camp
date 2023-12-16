package tag_repository

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/ports/repository"
	"github.com/awlsring/camp/internal/pkg/domain/tag"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/pkg/errors"
)

func (r *TagRepo) DeleteTagFromResource(ctx context.Context, key tag.TagKey, rid string) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Deleting tag %s from resource %s", key, rid)

	query := `DELETE FROM resource_tags WHERE resource_identifier = $1 AND tag_key = $2;`
	log.Debug().Msgf("Query: %s", query)

	log.Debug().Msg("Deleting tag from database")
	_, err := r.database.ExecContext(ctx, query, rid, key)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting tag from database")
		return errors.Wrap(repository.ErrInternalFailure, err.Error())
	}

	log.Debug().Msg("Successfully deleted tag from resource")
	return nil
}
