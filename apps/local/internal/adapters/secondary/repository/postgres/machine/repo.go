package machine_repository

import (
	"github.com/awlsring/camp/apps/local/internal/ports/repository"
	"github.com/awlsring/camp/internal/pkg/database"

	_ "github.com/lib/pq"
)

var _ repository.Machine = &MachineRepo{}

type MachineRepo struct {
	database database.Database
	tagRepo  repository.Tag
}

func New(db database.Database, tagRepo repository.Tag) (repository.Machine, error) {
	r := &MachineRepo{
		database: db,
		tagRepo:  tagRepo,
	}
	err := r.init()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *MachineRepo) init() error {
	return r.initTables()
}

func (r *MachineRepo) Close() error {
	return r.database.Close()
}
