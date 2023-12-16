package topic

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/power"
)

type PowerStateChange interface {
	SendStateChangeMessage(ctx context.Context, msg *power.StateChangeMessage) error
}
