package service

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
)

type Machine interface {
	RegisterMachine(ctx context.Context, id machine.Identifier, endpoint machine.MachineEndpoint, key machine.AgentKey, class machine.MachineClass, sys *machine.System, cpu *machine.Cpu, mem *machine.Memory, disks []*machine.Disk, nics []*machine.NetworkInterface, vols []*machine.Volume, ips []*machine.IpAddress) error
	DescribeMachine(ctx context.Context, id machine.Identifier) (*machine.Machine, error)
	ListMachines(ctx context.Context) ([]*machine.Machine, error)
	AddTagsToMachine(ctx context.Context, id machine.Identifier, tags []*tag.Tag) error
	RemoveTagFromMachine(ctx context.Context, id machine.Identifier, tag tag.TagKey) error
	AcknowledgeHeartbeat(ctx context.Context, id machine.Identifier) error
	UpdateStatus(ctx context.Context, id machine.Identifier, status machine.MachineStatus) error
	ReportSystemChange(ctx context.Context, id machine.Identifier, class machine.MachineClass, sys *machine.System, cpu *machine.Cpu, mem *machine.Memory, disks []*machine.Disk, nics []*machine.NetworkInterface, vols []*machine.Volume, ips []*machine.IpAddress) error
}
