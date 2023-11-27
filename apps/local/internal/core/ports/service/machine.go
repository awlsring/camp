package service

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

type Machine interface {
	RegisterMachine(ctx context.Context, id machine.InternalIdentifier, class machine.MachineClass, sys *machine.System, cpu *machine.Cpu, mem *machine.Memory, disks []*machine.Disk, nics []*machine.NetworkInterface, vols []*machine.Volume, ips []*machine.IpAddress) error
	DescribeMachine(ctx context.Context, id machine.Identifier) (*machine.Machine, error)
	ListMachines(ctx context.Context) ([]*machine.Machine, error)
	AcknowledgeHeartbeat(ctx context.Context, id machine.InternalIdentifier) error
	UpdateStatus(ctx context.Context, id machine.InternalIdentifier, status machine.MachineStatus) error
	ReportSystemChange(ctx context.Context, id machine.InternalIdentifier, class machine.MachineClass, sys *machine.System, cpu *machine.Cpu, mem *machine.Memory, disks []*machine.Disk, nics []*machine.NetworkInterface, vols []*machine.Volume, ips []*machine.IpAddress) error
}
