package topic

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
)

type Power interface {
	SendPowerChangeRequest(ctx context.Context, msg *power.PowerChangeRequestMessage) error
}
