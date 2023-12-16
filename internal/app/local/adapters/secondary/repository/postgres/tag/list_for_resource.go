package tag_repository

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/tag"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (r *TagRepo) ListForResource(ctx context.Context, rid string) ([]*tag.Tag, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Listing tags for resource %s", rid)

	query := `SELECT * FROM resource_tags WHERE resource_identifier = $1;`
	log.Debug().Msgf("Query: %s", query)

	log.Debug().Msg("Selecting tags from database")
	var tagModels []*TagSql
	err := r.database.SelectContext(ctx, &tagModels, query, rid)
	if err != nil {
		log.Error().Err(err).Msg("Error selecting tags from database")
		return nil, err
	}

	log.Debug().Msg("Converting tag models to tags")
	tags := make([]*tag.Tag, len(tagModels))
	for i, tagModel := range tagModels {
		tag, err := tagModel.ToModel()
		if err != nil {
			log.Error().Err(err).Msg("Error converting tag model to tag")
			return nil, err
		}
		tags[i] = tag
	}

	log.Debug().Msgf("Returning %d tags", len(tags))
	return tags, nil
}
