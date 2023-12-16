package service

import (
	"context"

	mach "github.com/awlsring/camp/internal/app/local/core/domain/machine"
	pwr "github.com/awlsring/camp/internal/app/local/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/domain/storage"
	"github.com/awlsring/camp/internal/pkg/domain/tag"
)

type Machine interface {
	RegisterMachine(ctx context.Context, id mach.Identifier, endpoint mach.MachineEndpoint, key mach.AgentKey, class machine.MachineClass, cap mach.PowerCapabilities, host *host.Host, cpu *cpu.CPU, mem *memory.Memory, disks []*storage.Disk, nics []*network.Nic, vols []*storage.Volume, ips []*network.IpAddress) error
	DescribeMachine(ctx context.Context, id mach.Identifier) (*mach.Machine, error)
	ListMachines(ctx context.Context) ([]*mach.Machine, error)
	RequestPowerChange(ctx context.Context, id mach.Identifier, changeType pwr.ChangeType) error
	AddTagsToMachine(ctx context.Context, id mach.Identifier, tags []*tag.Tag) error
	RemoveTagFromMachine(ctx context.Context, id mach.Identifier, tag tag.TagKey) error
	AcknowledgeHeartbeat(ctx context.Context, id mach.Identifier) error
	UpdateStatus(ctx context.Context, id mach.Identifier, status power.StatusCode) error
	ReportSystemChange(ctx context.Context, id mach.Identifier, host *host.Host, cpu *cpu.CPU, mem *memory.Memory, disks []*storage.Disk, nics []*network.Nic, vols []*storage.Volume, ips []*network.IpAddress) error
}
