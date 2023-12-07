package topic

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
)

type PowerStateJob interface {
	SendRequestStateChangeMessage(ctx context.Context, msg *power.RequestStateChangeMessage) error
	SendStateValidationMessage(ctx context.Context, msg *power.StateValidationMessage) error
}
