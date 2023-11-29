package tag_repository

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
	"github.com/awlsring/camp/apps/local/internal/ports/repository"
	"github.com/awlsring/camp/internal/pkg/database"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ repository.Tag = &TagRepo{}

type TagRepo struct {
	database database.Database
}

func New(db *sqlx.DB) (repository.Tag, error) {
	r := &TagRepo{
		database: db,
	}
	err := r.init()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *TagRepo) init() error {
	return r.initTable()
}

func (r *TagRepo) initTable() error {
	createTable := `
		CREATE TABLE IF NOT EXISTS resource_tags (
			id SERIAL PRIMARY KEY,
			resource_identifier VARCHAR(64),
			tag_key VARCHAR(50) NOT NULL,
			tag_value VARCHAR(128) NOT NULL,
			resource_type VARCHAR(50) NOT NULL
		);
	`

	tableQueries := []string{
		createTable,
	}

	ctx := context.Background()
	for _, query := range tableQueries {
		_, err := r.database.ExecContext(ctx, query)
		if err != nil {
			log.Error().Err(err).Msgf("Error creating table with query %s", query)
			return err
		}
	}
	return nil
}

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
		log.Error().Err(err).Msg("Error inserting tag into database")
		return err
	}

	log.Debug().Msg("Successfully added tag to resource")
	return nil
}

func (r *TagRepo) DeleteTagFromResource(ctx context.Context, key, rid string) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Deleting tag %s from resource %s", key, rid)

	query := `DELETE FROM resource_tags WHERE resource_identifier = $1 AND tag_key = $2;`
	log.Debug().Msgf("Query: %s", query)

	log.Debug().Msg("Deleting tag from database")
	_, err := r.database.ExecContext(ctx, query, rid, key)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting tag from database")
		return err
	}

	log.Debug().Msg("Successfully deleted tag from resource")
	return nil
}
