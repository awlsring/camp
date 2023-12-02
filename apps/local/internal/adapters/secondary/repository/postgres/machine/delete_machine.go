package machine_repository

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

func (r *MachineRepo) Delete(ctx context.Context, id machine.Identifier) error {
	_, err := r.database.ExecContext(ctx, "DELETE FROM machines WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
