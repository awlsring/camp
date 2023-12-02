package machine_repository

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

func (r MachineRepo) UpdateStatus(ctx context.Context, id machine.Identifier, status machine.MachineStatus) error {
	_, err := r.database.ExecContext(ctx, "UPDATE machines SET status = $1 WHERE identifier = $2", status.String(), id)
	if err != nil {
		return err
	}
	return nil
}
