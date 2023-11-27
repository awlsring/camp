package machine

import (
	"time"
)

type Machine struct {
	Identifier         Identifier
	InternalIdentifier InternalIdentifier
	Class              MachineClass
	LastHeartbeat      time.Time
	RegisteredAt       time.Time
	UpdatedAt          time.Time
	Status             MachineStatus
	System             *System
	Cpu                *Cpu
	Memory             *Memory
	Disks              []*Disk
	NetworkInterfaces  []*NetworkInterface
	Volumes            []*Volume
	Addresses          []*IpAddress
}

func NewMachine(identifier Identifier, internalId InternalIdentifier, class MachineClass, lastHeartbeat time.Time, registeredAt time.Time, updatedAt time.Time, status MachineStatus, system *System, cpu *Cpu, memory *Memory, disks []*Disk, networkInterfaces []*NetworkInterface, volumes []*Volume, addresses []*IpAddress) *Machine {
	return &Machine{
		Identifier:         identifier,
		InternalIdentifier: internalId,
		Class:              class,
		LastHeartbeat:      lastHeartbeat,
		RegisteredAt:       registeredAt,
		UpdatedAt:          updatedAt,
		Status:             status,
		System:             system,
		Cpu:                cpu,
		Memory:             memory,
		Disks:              disks,
		NetworkInterfaces:  networkInterfaces,
		Volumes:            volumes,
		Addresses:          addresses,
	}
}
