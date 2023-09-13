package machine

import "strings"

type MachineStatus string

const (
	Running              MachineStatus = "Running"
	Stopped              MachineStatus = "Stopped"
	Stopping             MachineStatus = "Stopping"
	Starting             MachineStatus = "Starting"
	Restarting           MachineStatus = "Restarting"
	MachineStatusUnknown MachineStatus = "Unknown"
)

type MachineClass string

const (
	BareMetal           MachineClass = "BareMetal"
	Virtual             MachineClass = "VirtualMachine"
	Hypervisor          MachineClass = "Hypervisor"
	MachineClassUnknown MachineClass = "Unknown"
)

type CpuArchitecture string

const (
	X86                    CpuArchitecture = "x86"
	Arm                    CpuArchitecture = "arm"
	CpuArchitectureUnknown CpuArchitecture = "Unknown"
)

type DiskInterface string

const (
	SATA                 DiskInterface = "SATA"
	SCSI                 DiskInterface = "SCSI"
	PCIe                 DiskInterface = "PCIe"
	DiskInterfaceUnknown DiskInterface = "Unknown"
)

type DiskType string

const (
	HDD             DiskType = "HDD"
	SSD             DiskType = "SSD"
	DiskTypeUnknown DiskType = "Unknown"
)

type IpAddressType string

const (
	V4                   IpAddressType = "v4"
	V6                   IpAddressType = "v6"
	IpAddressTypeUnknown IpAddressType = "Unknown"
)

func strToMachineClass(s string) MachineClass {
	switch strings.ToLower(s) {
	case "baremetal":
		return BareMetal
	case "virtualmachine":
		return Virtual
	case "hypervisor":
		return Hypervisor
	default:
		return MachineClassUnknown
	}
}

func strToCpuArchitecture(s string) CpuArchitecture {
	switch strings.ToLower(s) {
	case "x86":
		return X86
	case "arm":
		return Arm
	default:
		return CpuArchitectureUnknown
	}
}

func strToDiskInterface(s string) DiskInterface {
	switch strings.ToLower(s) {
	case "sata":
		return SATA
	case "scsi":
		return SCSI
	case "pcie":
		return PCIe
	default:
		return DiskInterfaceUnknown
	}
}

func strToDiskType(s string) DiskType {
	switch strings.ToLower(s) {
	case "hdd":
		return HDD
	case "ssd":
		return SSD
	default:
		return DiskTypeUnknown
	}
}

func strToIpAddressType(s string) IpAddressType {
	switch strings.ToLower(s) {
	case "v4":
		return V4
	case "v6":
		return V6
	default:
		return IpAddressTypeUnknown
	}
}
