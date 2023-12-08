package machine_repository

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (r MachineRepo) UpdateStatus(ctx context.Context, id machine.Identifier, status machine.MachineStatus) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Updating status to %s", status.String())

	_, err := r.database.ExecContext(ctx, "UPDATE power_state SET state = $1, updated_at = NOW() AT TIME ZONE 'UTC' WHERE machine_id = $2", status.String(), id)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to update status to %s", status.String())
		return err
	}

	log.Debug().Msgf("Status updated to %s", status.String())
	return nil
}
