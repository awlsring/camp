package topic

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
)

type PowerStateChange interface {
	SendStateChangeMessage(ctx context.Context, msg *power.StateChangeMessage) error
}
