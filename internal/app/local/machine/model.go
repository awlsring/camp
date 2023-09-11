package machine

type Model struct {
	Identifier        string
	Class             MachineClass
	System            *SystemModel
	Cpu               *CpuModel
	Memory            *MemoryModel
	Disks             []*DiskModel
	NetworkInterfaces []*NetworkInterfaceModel
	Volumes           []*VolumeModel
	Addresses         []*AddressModel
}

type SystemModel struct {
	Family        *string
	KernelVersion *string
	Os            *string
	OsVersion     *string
	OsPretty      *string
	Hostname      *string
}

type CpuModel struct {
	Cores        int
	Architecture CpuArchitecture
	Model        *string
	Vendor       *string
}

type MemoryModel struct {
	Total int64
}

type DiskModel struct {
	Device     string
	Model      *string
	Vendor     *string
	Interface  DiskInterface
	Type       DiskType
	Serial     *string
	SectorSize int
	Size       int64
	SizeRaw    *int64
}

type NetworkInterfaceModel struct {
	Name        string
	IpAddresses []AddressModel
	Virtual     bool
	MacAddress  *string
	Vendor      *string
	Mtu         *int16
	Speed       *int
	Duplex      *string
}

type VolumeModel struct {
	Name       string
	MountPoint string
	Total      int64
	FileSystem *string
}

type AddressModel struct {
	Version IpAddressType
	Address string
}
