package service

import (
	"context"
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

type PowerState interface {
	WakeOnLan(ctx context.Context, id machine.Identifier, mac machine.MacAddress) error
	PowerOff(ctx context.Context, id machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) error
	Reboot(ctx context.Context, id machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) error
	VerifyState(ctx context.Context, identifier machine.Identifier, state machine.MachineStatus, endpoint machine.MachineEndpoint, token machine.AgentKey) error
	VerifyTransitionalState(ctx context.Context, identifier machine.Identifier, started time.Time, deadline time.Time, state machine.MachineStatus, endpoint machine.MachineEndpoint, token machine.AgentKey) error
	ReconcileUnknownState(ctx context.Context, identifier machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) error
}
