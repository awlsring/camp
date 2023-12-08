package service

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

type PowerState interface {
	WakeOnLan(ctx context.Context, id machine.Identifier) error
	PowerOff(ctx context.Context, id machine.Identifier) error
	Reboot(ctx context.Context, id machine.Identifier) error
	VerifyFinalState(ctx context.Context, identifier machine.Identifier, state machine.MachineStatus) error
	VerifyTransitionalState(ctx context.Context, identifier machine.Identifier, state machine.MachineStatus) error
	ReconcileUnknownState(ctx context.Context, identifier machine.Identifier) error
}
