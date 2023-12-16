package machine_repository

import (
	"time"

	mach "github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/domain/storage"
)

type MachineSql struct {
	Identifier        string    `db:"identifier"`
	Class             string    `db:"class"`
	Endpoint          string    `db:"endpoint"`
	Key               string    `db:"key"`
	LastHeartbeat     time.Time `db:"last_heartbeat"`
	RegisteredAt      time.Time `db:"registered_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	Status            *StatusModelSql
	PowerCapabilities *PowerCapabilityModelSql
	Host              *HostModelSql
	Cpu               *CpuModelSql
	Memory            *MemoryModelSql
	Disks             []*DiskModelSql
	NetworkInterfaces []*NetworkInterfaceModelSql
	Volumes           []*VolumeModelSql
	Addresses         []*IpAddressModelSql
}

func (m *MachineSql) ToModel() (*mach.Machine, error) {
	host := m.Host.ToModel()
	cpu := m.Cpu.ToModel()
	memory := m.Memory.ToModel()

	disks := make([]*storage.Disk, len(m.Disks))
	for i, d := range m.Disks {
		di, err := d.ToModel()
		if err != nil {
			return nil, err
		}
		disks[i] = di
	}

	networkInterfaces := make([]*network.Nic, len(m.NetworkInterfaces))
	for i, n := range m.NetworkInterfaces {
		nic, err := n.ToModel()
		if err != nil {
			return nil, err
		}
		networkInterfaces[i] = nic
	}

	volumes := make([]*storage.Volume, len(m.Volumes))
	for i, v := range m.Volumes {
		vol, err := v.ToModel()
		if err != nil {
			return nil, err
		}
		volumes[i] = vol
	}

	addresses := make([]*network.IpAddress, len(m.Addresses))
	for i, a := range m.Addresses {
		add, err := a.ToModel()
		if err != nil {
			return nil, err
		}
		addresses[i] = add
	}

	id, err := mach.IdentifierFromString(m.Identifier)
	if err != nil {
		return nil, err
	}

	class, err := machine.MachineClassFromString(m.Class)
	if err != nil {
		return nil, err
	}

	endpoint, err := mach.MachineEndpointFromString(m.Endpoint)
	if err != nil {
		return nil, err
	}

	apiKey, err := mach.AgentKeyFromString(m.Key)
	if err != nil {
		return nil, err
	}

	return &mach.Machine{
		Identifier:        id,
		Class:             class,
		AgentEndpoint:     endpoint,
		AgentApiKey:       apiKey,
		LastHeartbeat:     m.LastHeartbeat,
		RegisteredAt:      m.RegisteredAt,
		UpdatedAt:         m.UpdatedAt,
		Status:            m.Status.ToModel(),
		PowerCapabilities: m.PowerCapabilities.ToModel(),
		Host:              host,
		Cpu:               cpu,
		Memory:            memory,
		Disks:             disks,
		NetworkInterfaces: networkInterfaces,
		Volumes:           volumes,
		Addresses:         addresses,
	}, nil

}

type StatusModelSql struct {
	State     string    `db:"state"`
	UpdatedAt time.Time `db:"updated_at"`
	Id        int64     `db:"id"`
	MachineId string    `db:"machine_id"`
}

func (p *StatusModelSql) ToModel() *power.Status {
	state, err := power.StatusCodeFromString(p.State)
	if err != nil {
		state = power.StatusCodeUnknown
	}

	return power.NewStatus(state, p.UpdatedAt)
}

type PowerCapabilityModelSql struct {
	Reboot       bool    `db:"reboot_enabled"`
	PowerOff     bool    `db:"power_off_enabled"`
	WakeOnLan    bool    `db:"wake_on_lan_enabled"`
	WakeOnLanMac *string `db:"wake_on_lan_mac,omitempty"`
	Id           int64   `db:"id"`
	MachineId    string  `db:"machine_id"`
}

func (p *PowerCapabilityModelSql) ToModel() mach.PowerCapabilities {
	var mac *network.MacAddress
	if p.WakeOnLanMac != nil {
		m, _ := network.MacAddressFromString(*p.WakeOnLanMac)
		mac = &m
	}

	return mach.PowerCapabilities{
		Reboot: mach.PowerCapabilityReboot{
			Enabled: p.Reboot,
		},
		PowerOff: mach.PowerCapabilityPowerOff{
			Enabled: p.PowerOff,
		},
		WakeOnLan: mach.PowerCapabilityWakeOnLan{
			Enabled:    p.WakeOnLan,
			MacAddress: mac,
		},
	}
}

type HostModelSql struct {
	OsPlatform *string `db:"platform,omitempty"`
	OsName     *string `db:"os_name,omitempty"`
	OsVersion  *string `db:"os_version,omitempty"`
	OsFamily   *string `db:"os_family,omitempty"`
	Kernel     *string `db:"kernel,omitempty"`
	Hostname   *string `db:"hostname,omitempty"`
	HostId     *string `db:"host_id,omitempty"`
	Id         int64   `db:"id"`
	MachineId  string  `db:"machine_id"`
}

func (s *HostModelSql) ToModel() *host.Host {
	return &host.Host{
		Hostname: s.Hostname,
		HostId:   s.HostId,
		OS: &host.OS{
			Family:   s.OsFamily,
			Kernel:   s.Kernel,
			Name:     s.OsName,
			Platform: s.OsPlatform,
			Version:  s.OsVersion,
		},
	}
}

type CpuModelSql struct {
	TotalCores   uint32  `db:"total_cores"`
	TotalThreads uint32  `db:"total_threads"`
	Architecture string  `db:"architecture"`
	Model        *string `db:"model,omitempty"`
	Vendor       *string `db:"vendor,omitempty"`
	Processors   []*ProcessorModelSql
	Id           int64  `db:"id"`
	MachineId    string `db:"machine_id"`
}

func (c *CpuModelSql) ToModel() *cpu.CPU {
	processors := make([]*cpu.Processor, len(c.Processors))
	for i, p := range c.Processors {
		processors[i] = p.ToModel()
	}
	return &cpu.CPU{
		TotalCores:   c.TotalCores,
		TotalThreads: c.TotalThreads,
		Architecture: cpu.ArchitectureFromString(c.Architecture),
		Model:        c.Model,
		Vendor:       c.Vendor,
		Processors:   processors,
	}
}

type ProcessorModelSql struct {
	Identifier  int     `db:"identifier"`
	CoreCount   uint32  `db:"core_count"`
	ThreadCount uint32  `db:"thread_count"`
	Model       *string `db:"model,omitempty"`
	Vendor      *string `db:"vendor,omitempty"`
	Id          int64   `db:"id"`
	CpuId       string  `db:"cpu_id"`
	Cores       []*CoreModelSql
}

func (p *ProcessorModelSql) ToModel() *cpu.Processor {
	cores := make([]*cpu.Core, len(p.Cores))
	for i, c := range p.Cores {
		cores[i] = c.ToModel()
	}
	return &cpu.Processor{
		Id:          p.Identifier,
		CoreCount:   p.CoreCount,
		ThreadCount: p.ThreadCount,
		Model:       p.Model,
		Vendor:      p.Vendor,
		Cores:       cores,
	}
}

type CoreModelSql struct {
	Identifier  int    `db:"identifier"`
	Threads     uint32 `db:"threads"`
	Id          int64  `db:"id"`
	ProcessorId int64  `db:"processor_id"`
}

func (c *CoreModelSql) ToModel() *cpu.Core {
	return &cpu.Core{
		Id:      c.Identifier,
		Threads: c.Threads,
	}
}

type MemoryModelSql struct {
	Total     uint64 `db:"total"`
	Id        int64  `db:"id"`
	MachineId string `db:"machine_id"`
}

func (m *MemoryModelSql) ToModel() *memory.Memory {
	return &memory.Memory{
		Total: m.Total,
	}
}

type DiskModelSql struct {
	Name              string  `db:"name"`
	Size              uint64  `db:"size"`
	DriveType         string  `db:"drive_type"`
	StorageController string  `db:"storage_controller"`
	Removable         bool    `db:"removable"`
	Model             *string `db:"model,omitempty"`
	Vendor            *string `db:"vendor,omitempty"`
	Serial            *string `db:"serial,omitempty"`
	WWN               *string `db:"wwn,omitempty"`
	Partitions        []*PartitionModelSql
	Id                int64  `db:"id"`
	MachineId         string `db:"machine_id"`
}

func (d *DiskModelSql) ToModel() (*storage.Disk, error) {
	driveType := storage.DriveTypeFromString(d.DriveType)
	controller := storage.StorageControllerFromString(d.StorageController)

	partitions := make([]*storage.Partition, len(d.Partitions))
	for i, p := range d.Partitions {
		part, err := p.ToModel()
		if err != nil {
			return nil, err
		}
		partitions[i] = part
	}

	return &storage.Disk{
		Name:              d.Name,
		Size:              d.Size,
		DriveType:         driveType,
		StorageController: controller,
		Removable:         d.Removable,
		Model:             d.Model,
		Vendor:            d.Vendor,
		Serial:            d.Serial,
		WWN:               d.WWN,
		Partitions:        partitions,
	}, nil
}

type PartitionModelSql struct {
	Name       string  `db:"name"`
	Size       uint64  `db:"size"`
	Readonly   bool    `db:"readonly"`
	Label      *string `db:"label,omitempty"`
	Type       *string `db:"type,omitempty"`
	FileSystem *string `db:"file_system,omitempty"`
	UUID       *string `db:"uuid,omitempty"`
	MountPoint *string `db:"mount_point,omitempty"`
	Id         int64   `db:"id"`
	DiskId     string  `db:"disk_id"`
}

func (p *PartitionModelSql) ToModel() (*storage.Partition, error) {
	return &storage.Partition{
		Name:       p.Name,
		Size:       p.Size,
		Readonly:   p.Readonly,
		Label:      p.Label,
		Type:       p.Type,
		FileSystem: p.FileSystem,
		UUID:       p.UUID,
		MountPoint: p.MountPoint,
	}, nil
}

type NetworkInterfaceModelSql struct {
	Name       string  `db:"name"`
	Virtual    bool    `db:"virtual"`
	MacAddress *string `db:"mac_address,omitempty"`
	Vendor     *string `db:"vendor,omitempty"`
	Speed      *string `db:"speed,omitempty"`
	Duplex     *string `db:"duplex,omitempty"`
	Id         int64   `db:"id"`
	MachineId  string  `db:"machine_id"`
	Addresses  []*IpAddressModelSql
}

func (n *NetworkInterfaceModelSql) ToModel() (*network.Nic, error) {
	addresses := make([]*network.IpAddress, len(n.Addresses))
	for i, a := range n.Addresses {
		add, err := a.ToModel()
		if err != nil {
			return nil, err
		}
		addresses[i] = add
	}

	var mac *network.MacAddress
	if n.MacAddress == nil {
		mac = nil
	} else if *n.MacAddress == "" {
		mac = nil
	} else {
		m, err := network.MacAddressFromString(*n.MacAddress)
		if err != nil {
			return nil, err
		}
		mac = &m
	}

	return &network.Nic{
		Name:        n.Name,
		Virtual:     n.Virtual,
		MacAddress:  mac,
		Vendor:      n.Vendor,
		Speed:       n.Speed,
		Duplex:      n.Duplex,
		IpAddresses: addresses,
	}, nil
}

type VolumeModelSql struct {
	Name       string  `db:"name"`
	MountPoint string  `db:"mount_point"`
	Total      uint64  `db:"total"`
	FileSystem *string `db:"file_system,omitempty"`
	Id         int64   `db:"id"`
	MachineId  string  `db:"machine_id"`
}

func (v *VolumeModelSql) ToModel() (*storage.Volume, error) {
	return &storage.Volume{
		Name:       v.Name,
		MountPoint: v.MountPoint,
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

func (a *IpAddressModelSql) ToModel() (*network.IpAddress, error) {
	ver := network.IpAddressTypeFromString(a.Version)

	add, err := network.AddressFromString(a.Address)
	if err != nil {
		return nil, err
	}

	return &network.IpAddress{
		Version: ver,
		Address: add,
	}, nil
}
