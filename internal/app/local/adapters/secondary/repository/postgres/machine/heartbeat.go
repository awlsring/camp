package machine_repository

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
)

func (r *Repository) UpdateHeartbeat(ctx context.Context, id machine.Identifier) error {
	_, err := r.database.ExecContext(ctx, "UPDATE machines SET last_heartbeat = NOW() WHERE identifier = $1", id)
	if err != nil {
		return err
	}
	return nil
}
