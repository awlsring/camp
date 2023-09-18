package machine

import "time"

type Model struct {
	Identifier        string                   `db:"identifier"`
	Class             MachineClass             `db:"class"`
	LastHeartbeat     time.Time                `db:"last_heartbeat"`
	RegisteredAt      time.Time                `db:"registered_at"`
	UpdatedAt         time.Time                `db:"updated_at"`
	Status            MachineStatus            `db:"status"`
	System            *SystemModel             `db:"system"`
	Cpu               *CpuModel                `db:"cpu"`
	Memory            *MemoryModel             `db:"memory"`
	Disks             []*DiskModel             `db:"disks"`
	NetworkInterfaces []*NetworkInterfaceModel `db:"network_interfaces"`
	Volumes           []*VolumeModel           `db:"volumes"`
	Addresses         []*AddressModel
}

type SystemModel struct {
	Family        *string `db:"family"`
	KernelVersion *string `db:"kernel_version"`
	Os            *string `db:"os"`
	OsVersion     *string `db:"os_version"`
	OsPretty      *string `db:"os_pretty"`
	Hostname      *string `db:"hostname"`
}

type CpuModel struct {
	Cores        int             `db:"cores"`
	Architecture CpuArchitecture `db:"architecture"`
	Model        *string         `db:"model"`
	Vendor       *string         `db:"vendor"`
}

type MemoryModel struct {
	Total int64 `db:"total"`
}

type DiskModel struct {
	Device     string        `db:"device"`
	Model      *string       `db:"model"`
	Vendor     *string       `db:"vendor"`
	Interface  DiskInterface `db:"interface"`
	Type       DiskType      `db:"type"`
	Serial     *string       `db:"serial"`
	SectorSize int           `db:"sector_size"`
	Size       int64         `db:"size"`
	SizeRaw    *int64        `db:"size_raw"`
}

type NetworkInterfaceModel struct {
	id          int64           `db:"id"`
	Name        string          `db:"name"`
	IpAddresses []*AddressModel `db:"ip_addresses"`
	Virtual     bool            `db:"virtual"`
	MacAddress  *string         `db:"mac_address"`
	Vendor      *string         `db:"vendor"`
	Mtu         *int            `db:"mtu"`
	Speed       *int            `db:"speed"`
	Duplex      *string         `db:"duplex"`
}

type VolumeModel struct {
	Name       string  `db:"name"`
	MountPoint string  `db:"mount_point"`
	Total      int64   `db:"total"`
	FileSystem *string `db:"file_system"`
}

type AddressModel struct {
	Version IpAddressType `db:"version"`
	Address string        `db:"address"`
}
