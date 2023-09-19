package repo

import (
	"time"

	"github.com/awlsring/camp/apps/local/machine"
)

type MachineSql struct {
	Identifier        string    `db:"identifier"`
	Class             string    `db:"class"`
	LastHeartbeat     time.Time `db:"last_heartbeat"`
	RegisteredAt      time.Time `db:"registered_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	Status            string    `db:"status"`
	System            *SystemModelSql
	Cpu               *CpuModelSql
	Memory            *MemoryModelSql
	Disks             []*DiskModelSql
	NetworkInterfaces []*NetworkInterfaceModelSql
	Volumes           []*VolumeModelSql
	Addresses         []*AddressModelSql
}

func (m *MachineSql) ToModel() *machine.Model {
	system := m.System.ToModel()
	cpu := m.Cpu.ToModel()
	memory := m.Memory.ToModel()

	disks := make([]*machine.DiskModel, len(m.Disks))
	for i, d := range m.Disks {
		disks[i] = d.ToModel()
	}

	networkInterfaces := make([]*machine.NetworkInterfaceModel, len(m.NetworkInterfaces))
	for i, n := range m.NetworkInterfaces {
		networkInterfaces[i] = n.ToModel()
	}

	volumes := make([]*machine.VolumeModel, len(m.Volumes))
	for i, v := range m.Volumes {
		volumes[i] = v.ToModel()
	}

	addresses := make([]*machine.AddressModel, len(m.Addresses))
	for i, a := range m.Addresses {
		addresses[i] = a.ToModel()
	}

	return &machine.Model{
		Identifier:        m.Identifier,
		Class:             machine.MachineClassFromString(m.Class),
		LastHeartbeat:     m.LastHeartbeat,
		RegisteredAt:      m.RegisteredAt,
		UpdatedAt:         m.UpdatedAt,
		Status:            machine.MachineStatusFromString(m.Status),
		System:            system,
		Cpu:               cpu,
		Memory:            memory,
		Disks:             disks,
		NetworkInterfaces: networkInterfaces,
		Volumes:           volumes,
		Addresses:         addresses,
	}

}

type SystemModelSql struct {
	Family        *string `db:"family"`
	KernelVersion *string `db:"kernel_version"`
	Os            *string `db:"os"`
	OsVersion     *string `db:"os_version"`
	OsPretty      *string `db:"os_pretty"`
	Hostname      *string `db:"hostname"`
	Id            int64   `db:"id"`
	MachineId     string  `db:"machine_id"`
}

func (s *SystemModelSql) ToModel() *machine.SystemModel {
	return &machine.SystemModel{
		Family:        s.Family,
		KernelVersion: s.KernelVersion,
		Os:            s.Os,
		OsVersion:     s.OsVersion,
		OsPretty:      s.OsPretty,
		Hostname:      s.Hostname,
	}
}

type CpuModelSql struct {
	Cores        int     `db:"cores"`
	Architecture string  `db:"architecture"`
	Model        *string `db:"model"`
	Vendor       *string `db:"vendor"`
	Id           int64   `db:"id"`
	MachineId    string  `db:"machine_id"`
}

func (c *CpuModelSql) ToModel() *machine.CpuModel {
	return &machine.CpuModel{
		Cores:        c.Cores,
		Architecture: machine.CpuArchitectureFromString(c.Architecture),
		Model:        c.Model,
		Vendor:       c.Vendor,
	}
}

type MemoryModelSql struct {
	Total     int64  `db:"total"`
	Id        int64  `db:"id"`
	MachineId string `db:"machine_id"`
}

func (m *MemoryModelSql) ToModel() *machine.MemoryModel {
	return &machine.MemoryModel{
		Total: m.Total,
	}
}

type DiskModelSql struct {
	Device     string  `db:"device"`
	Model      *string `db:"model"`
	Vendor     *string `db:"vendor"`
	Interface  string  `db:"interface"`
	Type       string  `db:"type"`
	Serial     *string `db:"serial"`
	SectorSize int     `db:"sector_size"`
	Size       int64   `db:"size"`
	SizeRaw    *int64  `db:"size_raw"`
	Id         int64   `db:"id"`
	MachineId  string  `db:"machine_id"`
}

func (d *DiskModelSql) ToModel() *machine.DiskModel {
	return &machine.DiskModel{
		Device:     d.Device,
		Model:      d.Model,
		Vendor:     d.Vendor,
		Interface:  machine.DiskInterfaceFromString(d.Interface),
		Type:       machine.DiskTypeFromString(d.Type),
		Serial:     d.Serial,
		SectorSize: d.SectorSize,
		Size:       d.Size,
		SizeRaw:    d.SizeRaw,
	}
}

type NetworkInterfaceModelSql struct {
	Name       string  `db:"name"`
	Virtual    bool    `db:"virtual"`
	MacAddress *string `db:"mac_address"`
	Vendor     *string `db:"vendor"`
	Mtu        *int    `db:"mtu"`
	Speed      *int    `db:"speed"`
	Duplex     *string `db:"duplex"`
	Id         int64   `db:"id"`
	MachineId  string  `db:"machine_id"`
	Addresses  []*AddressModelSql
}

func (n *NetworkInterfaceModelSql) ToModel() *machine.NetworkInterfaceModel {
	addresses := make([]*machine.AddressModel, len(n.Addresses))
	for i, a := range n.Addresses {
		addresses[i] = a.ToModel()
	}

	return &machine.NetworkInterfaceModel{
		Name:        n.Name,
		Virtual:     n.Virtual,
		MacAddress:  n.MacAddress,
		Vendor:      n.Vendor,
		Mtu:         n.Mtu,
		Speed:       n.Speed,
		Duplex:      n.Duplex,
		IpAddresses: addresses,
	}
}

type VolumeModelSql struct {
	Name       string  `db:"name"`
	MountPoint string  `db:"mount_point"`
	Total      int64   `db:"total"`
	FileSystem *string `db:"file_system"`
	Id         int64   `db:"id"`
	MachineId  string  `db:"machine_id"`
}

func (v *VolumeModelSql) ToModel() *machine.VolumeModel {
	return &machine.VolumeModel{
		Name:       v.Name,
		MountPoint: v.MountPoint,
		Total:      v.Total,
		FileSystem: v.FileSystem,
	}
}

type AddressModelSql struct {
	Version string `db:"version"`
	Address string `db:"address"`
	Id      int64  `db:"id"`
	NicId   string `db:"nic_id"`
}

func (a *AddressModelSql) ToModel() *machine.AddressModel {
	return &machine.AddressModel{
		Version: machine.IpAddressTypeFromString(a.Version),
		Address: a.Address,
	}
}
