package tag_repository

import (
	"context"
	"strings"

	"github.com/awlsring/camp/internal/app/local/ports/repository"
	"github.com/awlsring/camp/internal/pkg/domain/tag"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (r *TagRepo) AddToResource(ctx context.Context, t *tag.Tag, id string, typee tag.ResourceType) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Adding tag %s to resource %s", t.Key, id)

	query := `
		INSERT INTO resource_tags (resource_identifier, tag_key, tag_value, resource_type)
		VALUES ($1, $2, $3, $4);`
	log.Debug().Msgf("Query: %s", query)

	log.Debug().Msg("Inserting tag into database")
	_, err := r.database.ExecContext(ctx, query, id, t.Key, t.Value, typee.String())
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			log.Warn().Err(err).Msg("Duplicate tag for resource")
			return repository.ErrDuplicateResourceTag
		}
		log.Error().Err(err).Msg("Error inserting tag into database")
		return err
	}

	log.Debug().Msg("Successfully added tag to resource")
	return nil
}
