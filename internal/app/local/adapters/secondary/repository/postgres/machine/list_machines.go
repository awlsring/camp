package machine_repository

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/app/local/ports/repository"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (r *Repository) List(ctx context.Context, filters *repository.ListMachinesFilters) ([]*machine.Machine, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke Repository.GetMachines")

	log.Debug().Msg("Listing machines from database")
	var machinesModels []*MachineSql
	err := r.database.SelectContext(ctx, &machinesModels, "SELECT * FROM machines")
	if err != nil {
		log.Error().Err(err).Msg("Failed to list machines")
		return nil, err
	}

	log.Debug().Msgf("Found %d machines", len(machinesModels))
	models := []*machine.Machine{}
	for _, m := range machinesModels {
		log.Debug().Msgf("Enriching machine: %+v", m.Identifier)
		err = r.enrichMachineEntry(ctx, m)
		if err != nil {
			return nil, err
		}

		log.Debug().Interface("machine", m).Msg("Machine data")
		log.Debug().Msgf("Converting machine to domain: %+v", m.Identifier)
		mod, err := m.ToModel()
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert machine to domain")
			return nil, err
		}
		tags, err := r.tagRepo.ListForResource(ctx, m.Identifier)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list tags for machine")
			return nil, err
		}
		mod.Tags = tags

		log.Debug().Msgf("Appending machine to list: %+v", m.Identifier)
		models = append(models, mod)
	}

	log.Debug().Msgf("Returning %d machines", len(models))
	return models, nil
}
