package service

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

type PowerControl interface {
	WakeOnLan(ctx context.Context, id machine.Identifier) error
	PowerOff(ctx context.Context, id machine.Identifier) error
	Reboot(ctx context.Context, id machine.Identifier) error
}
