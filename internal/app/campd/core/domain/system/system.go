package system

import (
	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/domain/motherboard"
	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/domain/storage"
)

type System struct {
	Host              *host.Host
	Bios              *motherboard.Bios
	Motherboard       *motherboard.Motherboard
	Cpu               *cpu.CPU
	Memory            *memory.Memory
	Disks             []*storage.Disk
	NetworkInterfaces []*network.Nic
	IpAddresses       []*network.IpAddress
}

func NewSystem(host *host.Host, bios *motherboard.Bios, motherboard *motherboard.Motherboard, cpu *cpu.CPU, memory *memory.Memory, disks []*storage.Disk, networkInterfaces []*network.Nic, ipAddresses []*network.IpAddress) *System {
	return &System{
		Host:              host,
		Bios:              bios,
		Motherboard:       motherboard,
		Cpu:               cpu,
		Memory:            memory,
		Disks:             disks,
		NetworkInterfaces: networkInterfaces,
		IpAddresses:       ipAddresses,
	}
}
