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

func MachineStatusFromString(v string) MachineStatus {
	switch strings.ToLower(v) {
	case "running":
		return Running
	case "stopped":
		return Stopped
	case "stopping":
		return Stopping
	case "starting":
		return Starting
	case "restarting":
		return Restarting
	default:
		return MachineStatusUnknown
	}
}

type MachineClass string

const (
	BareMetal           MachineClass = "BareMetal"
	Virtual             MachineClass = "VirtualMachine"
	Hypervisor          MachineClass = "Hypervisor"
	MachineClassUnknown MachineClass = "Unknown"
)

func MachineClassFromString(v string) MachineClass {
	switch strings.ToLower(v) {
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

type CpuArchitecture string

const (
	X86                    CpuArchitecture = "x86"
	Arm                    CpuArchitecture = "arm"
	CpuArchitectureUnknown CpuArchitecture = "Unknown"
)

func CpuArchitectureFromString(v string) CpuArchitecture {
	switch strings.ToLower(v) {
	case "x86", "x86_64":
		return X86
	case "arm", "armv7", "armv8":
		return Arm
	default:
		return CpuArchitectureUnknown
	}
}

type DiskInterface string

const (
	SATA                 DiskInterface = "SATA"
	SCSI                 DiskInterface = "SCSI"
	PCIe                 DiskInterface = "PCIe"
	DiskInterfaceUnknown DiskInterface = "Unknown"
)

func DiskInterfaceFromString(v string) DiskInterface {
	switch strings.ToLower(v) {
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

type DiskType string

const (
	HDD             DiskType = "HDD"
	SSD             DiskType = "SSD"
	DiskTypeUnknown DiskType = "Unknown"
)

func DiskTypeFromString(v string) DiskType {
	switch strings.ToLower(v) {
	case "hdd":
		return HDD
	case "ssd":
		return SSD
	default:
		return DiskTypeUnknown
	}
}

type IpAddressType string

const (
	V4                   IpAddressType = "v4"
	V6                   IpAddressType = "v6"
	IpAddressTypeUnknown IpAddressType = "Unknown"
)

func IpAddressTypeFromString(v string) IpAddressType {
	switch strings.ToLower(v) {
	case "ipv4", "v4":
		return V4
	case "ipv6", "v6":
		return V6
	default:
		return IpAddressTypeUnknown
	}
}
