package service

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/power"
)

type PowerState interface {
	WakeOnLan(ctx context.Context, id machine.Identifier) error
	PowerOff(ctx context.Context, id machine.Identifier) error
	Reboot(ctx context.Context, id machine.Identifier) error
	VerifyFinalState(ctx context.Context, identifier machine.Identifier, state power.StatusCode) error
	VerifyTransitionalState(ctx context.Context, identifier machine.Identifier, state power.StatusCode) error
	ReconcileUnknownState(ctx context.Context, identifier machine.Identifier) error
}
