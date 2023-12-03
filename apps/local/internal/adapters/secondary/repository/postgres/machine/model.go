package machine_repository

import (
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

type MachineSql struct {
	Identifier        string    `db:"identifier"`
	Class             string    `db:"class"`
	LastHeartbeat     time.Time `db:"last_heartbeat"`
	RegisteredAt      time.Time `db:"registered_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	Status            string    `db:"status"`
	PowerCapabilities *PowerCapabilityModelSql
	System            *SystemModelSql
	Cpu               *CpuModelSql
	Memory            *MemoryModelSql
	Disks             []*DiskModelSql
	NetworkInterfaces []*NetworkInterfaceModelSql
	Volumes           []*VolumeModelSql
	Addresses         []*IpAddressModelSql
}

func (m *MachineSql) ToModel() (*machine.Machine, error) {
	system := m.System.ToModel()
	cpu := m.Cpu.ToModel()
	memory := m.Memory.ToModel()

	disks := make([]*machine.Disk, len(m.Disks))
	for i, d := range m.Disks {
		di, err := d.ToModel()
		if err != nil {
			return nil, err
		}
		disks[i] = di
	}

	networkInterfaces := make([]*machine.NetworkInterface, len(m.NetworkInterfaces))
	for i, n := range m.NetworkInterfaces {
		nic, err := n.ToModel()
		if err != nil {
			return nil, err
		}
		networkInterfaces[i] = nic
	}

	volumes := make([]*machine.Volume, len(m.Volumes))
	for i, v := range m.Volumes {
		vol, err := v.ToModel()
		if err != nil {
			return nil, err
		}
		volumes[i] = vol
	}

	addresses := make([]*machine.IpAddress, len(m.Addresses))
	for i, a := range m.Addresses {
		add, err := a.ToModel()
		if err != nil {
			return nil, err
		}
		addresses[i] = add
	}

	id, err := machine.IdentifierFromString(m.Identifier)
	if err != nil {
		return nil, err
	}

	class, err := machine.MachineClassFromString(m.Class)
	if err != nil {
		return nil, err
	}

	status, err := machine.MachineStatusFromString(m.Status)
	if err != nil {
		return nil, err
	}

	return &machine.Machine{
		Identifier:        id,
		Class:             class,
		LastHeartbeat:     m.LastHeartbeat,
		RegisteredAt:      m.RegisteredAt,
		UpdatedAt:         m.UpdatedAt,
		Status:            status,
		System:            system,
		Cpu:               cpu,
		Memory:            memory,
		Disks:             disks,
		NetworkInterfaces: networkInterfaces,
		Volumes:           volumes,
		Addresses:         addresses,
	}, nil

}

type PowerCapabilityModelSql struct {
	Reboot       bool    `db:"reboot_enabled"`
	PowerOff     bool    `db:"power_off_enabled"`
	WakeOnLan    bool    `db:"wake_on_lan_enabled"`
	WakeOnLanMac *string `db:"wake_on_lan_mac,omitempty"`
	Id           int64   `db:"id"`
	MachineId    string  `db:"machine_id"`
}

func (p *PowerCapabilityModelSql) ToModel() machine.PowerCapabilities {
	var mac *machine.MacAddress
	if p.WakeOnLanMac != nil {
		m, _ := machine.MacAddressFromString(*p.WakeOnLanMac)
		mac = &m
	}

	return machine.PowerCapabilities{
		Reboot: machine.PowerCapabilityReboot{
			Enabled: p.Reboot,
		},
		PowerOff: machine.PowerCapabilityPowerOff{
			Enabled: p.PowerOff,
		},
		WakeOnLan: machine.PowerCapabilityWakeOnLan{
			Enabled:    p.WakeOnLan,
			MacAddress: mac,
		},
	}
}

type SystemModelSql struct {
	Family        *string `db:"family,omitempty"`
	KernelVersion *string `db:"kernel_version,omitempty"`
	Os            *string `db:"os,omitempty"`
	OsVersion     *string `db:"os_version,omitempty"`
	OsPretty      *string `db:"os_pretty,omitempty"`
	Hostname      *string `db:"hostname,omitempty"`
	Id            int64   `db:"id"`
	MachineId     string  `db:"machine_id"`
}

func (s *SystemModelSql) ToModel() *machine.System {
	return &machine.System{
		Os: machine.Os{
			Family:     s.Family,
			Version:    s.OsVersion,
			Name:       s.Os,
			PrettyName: s.OsPretty,
			Kernel:     s.KernelVersion,
		},
		Hostname: s.Hostname,
	}
}

type CpuModelSql struct {
	Cores        int     `db:"cores"`
	Architecture string  `db:"architecture"`
	Model        *string `db:"model,omitempty"`
	Vendor       *string `db:"vendor,omitempty"`
	Id           int64   `db:"id"`
	MachineId    string  `db:"machine_id"`
}

func (c *CpuModelSql) ToModel() *machine.Cpu {
	return &machine.Cpu{
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

func (m *MemoryModelSql) ToModel() *machine.Memory {
	return &machine.Memory{
		Total: m.Total,
	}
}

type DiskModelSql struct {
	Device     string  `db:"device"`
	Model      *string `db:"model,omitempty"`
	Vendor     *string `db:"vendor,omitempty"`
	Interface  string  `db:"interface"`
	Type       string  `db:"type"`
	Serial     *string `db:"serial,omitempty"`
	SectorSize int     `db:"sector_size"`
	Size       int64   `db:"size"`
	SizeRaw    *int64  `db:"size_raw,omitempty"`
	Id         int64   `db:"id"`
	MachineId  string  `db:"machine_id"`
}

func (d *DiskModelSql) ToModel() (*machine.Disk, error) {
	dev, err := machine.DiskIdentifierFromString(d.Device)
	if err != nil {
		return nil, err
	}

	return &machine.Disk{
		Device:     dev,
		Model:      d.Model,
		Vendor:     d.Vendor,
		Interface:  machine.DiskInterfaceFromString(d.Interface),
		Type:       machine.DiskClassFromString(d.Type),
		Serial:     d.Serial,
		SectorSize: d.SectorSize,
		Size:       d.Size,
		SizeRaw:    d.SizeRaw,
	}, nil
}

type NetworkInterfaceModelSql struct {
	Name       string  `db:"name"`
	Virtual    bool    `db:"virtual"`
	MacAddress *string `db:"mac_address,omitempty"`
	Vendor     *string `db:"vendor,omitempty"`
	Mtu        *int    `db:"mtu,omitempty"`
	Speed      *int    `db:"speed,omitempty"`
	Duplex     *string `db:"duplex,omitempty"`
	Id         int64   `db:"id"`
	MachineId  string  `db:"machine_id"`
	Addresses  []*IpAddressModelSql
}

func (n *NetworkInterfaceModelSql) ToModel() (*machine.NetworkInterface, error) {
	addresses := make([]*machine.IpAddress, len(n.Addresses))
	for i, a := range n.Addresses {
		add, err := a.ToModel()
		if err != nil {
			return nil, err
		}
		addresses[i] = add
	}

	name, err := machine.NetworkInterfaceIdentifierFromString(n.Name)
	if err != nil {
		return nil, err
	}

	var mac *machine.MacAddress
	if n.MacAddress == nil {
		mac = nil
	} else if *n.MacAddress == "" {
		mac = nil
	} else {
		m, err := machine.MacAddressFromString(*n.MacAddress)
		if err != nil {
			return nil, err
		}
		mac = &m
	}

	return &machine.NetworkInterface{
		Name:        name,
		Virtual:     n.Virtual,
		MacAddress:  mac,
		Vendor:      n.Vendor,
		Mtu:         n.Mtu,
		Speed:       n.Speed,
		Duplex:      n.Duplex,
		IpAddresses: addresses,
	}, nil
}

type VolumeModelSql struct {
	Name       string  `db:"name"`
	MountPoint string  `db:"mount_point"`
	Total      int64   `db:"total"`
	FileSystem *string `db:"file_system,omitempty"`
	Id         int64   `db:"id"`
	MachineId  string  `db:"machine_id"`
}

func (v *VolumeModelSql) ToModel() (*machine.Volume, error) {
	vol, err := machine.VolumeIdentifierFromString(v.Name)
	if err != nil {
		return nil, err
	}

	m, err := machine.MountPointFromString(v.Name)
	if err != nil {
		return nil, err
	}

	return &machine.Volume{
		Name:       vol,
		MountPoint: m,
		Total:      v.Total,
		FileSystem: v.FileSystem,
	}, nil
}

type IpAddressModelSql struct {
	Version string `db:"version"`
	Address string `db:"address"`
	Id      int64  `db:"id"`
	NicId   string `db:"nic_id"`
}

func (a *IpAddressModelSql) ToModel() (*machine.IpAddress, error) {
	ver := machine.IpAddressTypeFromString(a.Version)

	add, err := machine.AddressFromString(a.Address)
	if err != nil {
		return nil, err
	}

	return &machine.IpAddress{
		Version: ver,
		Address: add,
	}, nil
}
