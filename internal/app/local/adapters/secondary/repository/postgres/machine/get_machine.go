package machine_repository

import (
	"context"
	"strings"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/app/local/ports/repository"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (r *Repository) Get(ctx context.Context, id machine.Identifier) (*machine.Machine, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Getting machine with id %s", id.String())

	var mdb MachineSql
	query := "SELECT * FROM machines WHERE identifier = $1"
	log.Debug().Msgf("Query: %s", query)
	err := r.database.GetContext(ctx, &mdb, query, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Warn().Msgf("Machine with id %s not found", id.String())
			return nil, repository.ErrMachineNotFound
		}
		log.Error().Err(err).Msgf("Failed to get machine with id %s", id.String())
		return nil, err
	}

	log.Debug().Msgf("Enriching machine: %+v", mdb.Identifier)
	err = r.enrichMachineEntry(ctx, &mdb)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Converting machine to domain: %+v", mdb.Identifier)
	mod, err := mdb.ToModel()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to convert machine to domain: %+v", mdb.Identifier)
		return nil, err
	}

	log.Debug().Msgf("Getting tags for machine: %+v", mdb.Identifier)
	tags, err := r.tagRepo.ListForResource(ctx, id.String())
	if err != nil {
		log.Error().Err(err).Msgf("Failed to list tags for machine: %+v", mdb.Identifier)
		return nil, err
	}
	log.Debug().Msgf("Found %d tags for machine: %+v", len(tags), mdb.Identifier)
	mod.Tags = tags

	log.Debug().Msgf("Returning machine: %+v", mdb.Identifier)
	return mod, nil
}
