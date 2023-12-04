package service

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

type PowerControl interface {
	WakeOnLan(ctx context.Context, id machine.Identifier, mac machine.MacAddress) error
	PowerOff(ctx context.Context, id machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) error
	Reboot(ctx context.Context, id machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) error
	CheckMachinePower(ctx context.Context, identifier machine.Identifier, endpoint machine.MachineEndpoint, token machine.AgentKey) (bool, error)
}
