package machine

import (
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
)

type Machine struct {
	Identifier        Identifier
	Class             MachineClass
	Tags              []*tag.Tag
	LastHeartbeat     time.Time
	RegisteredAt      time.Time
	UpdatedAt         time.Time
	Status            MachineStatus
	System            *System
	Cpu               *Cpu
	Memory            *Memory
	Disks             []*Disk
	NetworkInterfaces []*NetworkInterface
	Volumes           []*Volume
	Addresses         []*IpAddress
}

func NewMachine(identifier Identifier, class MachineClass, lastHeartbeat time.Time, registeredAt time.Time, updatedAt time.Time, status MachineStatus, system *System, cpu *Cpu, memory *Memory, disks []*Disk, networkInterfaces []*NetworkInterface, volumes []*Volume, addresses []*IpAddress) *Machine {
	return &Machine{
		Identifier:        identifier,
		Class:             class,
		LastHeartbeat:     lastHeartbeat,
		RegisteredAt:      registeredAt,
		UpdatedAt:         updatedAt,
		Status:            status,
		System:            system,
		Cpu:               cpu,
		Memory:            memory,
		Disks:             disks,
		NetworkInterfaces: networkInterfaces,
		Volumes:           volumes,
		Addresses:         addresses,
	}
}
