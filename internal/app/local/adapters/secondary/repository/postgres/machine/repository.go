package machine_repository

import (
	"github.com/awlsring/camp/internal/app/local/ports/repository"
	"github.com/awlsring/camp/internal/pkg/database"

	_ "github.com/lib/pq"
)

var _ repository.Machine = &Repository{}

type Repository struct {
	database database.Database
	tagRepo  repository.Tag
}

func New(db database.Database, tagRepo repository.Tag) (repository.Machine, error) {
	r := &Repository{
		database: db,
		tagRepo:  tagRepo,
	}
	err := r.init()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Repository) init() error {
	return r.initTables()
}

func (r *Repository) Close() error {
	return r.database.Close()
}
