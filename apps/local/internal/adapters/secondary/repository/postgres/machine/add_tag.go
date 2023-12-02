package machine_repository

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (r *MachineRepo) AddTags(ctx context.Context, id machine.Identifier, tags []*tag.Tag) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Adding tags to machine %s", id.String())

	for _, t := range tags {
		err := r.tagRepo.AddToResource(ctx, t, id.String(), tag.ResourceTypeMachine)
		if err != nil {
			return err
		}
	}

	return nil
}
