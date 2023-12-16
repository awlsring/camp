package machine

import (
	"time"

	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/domain/storage"
	"github.com/awlsring/camp/internal/pkg/domain/tag"
)

type Machine struct {
	Identifier        Identifier
	AgentEndpoint     MachineEndpoint
	AgentApiKey       AgentKey
	PowerCapabilities PowerCapabilities
	Class             machine.MachineClass
	Tags              []*tag.Tag
	LastHeartbeat     time.Time
	RegisteredAt      time.Time
	UpdatedAt         time.Time
	Status            *power.Status
	Host              *host.Host
	Cpu               *cpu.CPU
	Memory            *memory.Memory
	Disks             []*storage.Disk
	NetworkInterfaces []*network.Nic
	Volumes           []*storage.Volume
	Addresses         []*network.IpAddress
}

func NewMachine(identifier Identifier, endpoint MachineEndpoint, key AgentKey, class machine.MachineClass, lastHeartbeat time.Time, registeredAt time.Time, updatedAt time.Time, status *power.Status, cap PowerCapabilities, host *host.Host, cpu *cpu.CPU, memory *memory.Memory, disks []*storage.Disk, networkInterfaces []*network.Nic, volumes []*storage.Volume, addresses []*network.IpAddress) *Machine {
	return &Machine{
		Identifier:        identifier,
		AgentEndpoint:     endpoint,
		AgentApiKey:       key,
		Class:             class,
		LastHeartbeat:     lastHeartbeat,
		RegisteredAt:      registeredAt,
		UpdatedAt:         updatedAt,
		Status:            status,
		PowerCapabilities: cap,
		Host:              host,
		Cpu:               cpu,
		Memory:            memory,
		Disks:             disks,
		NetworkInterfaces: networkInterfaces,
		Volumes:           volumes,
		Addresses:         addresses,
	}
}
