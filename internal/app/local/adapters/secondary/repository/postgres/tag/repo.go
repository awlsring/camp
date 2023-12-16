package tag_repository

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/ports/repository"
	"github.com/awlsring/camp/internal/pkg/database"
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
			resource_type VARCHAR(50) NOT NULL,
			CONSTRAINT unique_resource_tag_key UNIQUE (resource_identifier, tag_key)
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
