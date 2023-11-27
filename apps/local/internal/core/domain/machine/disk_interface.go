package machine

import (
	"strings"
)

type DiskInterface int64

const (
	DiskInterfaceSata DiskInterface = iota
	DiskInterfaceScsi
	DiskInterfaceSas
	DiskInterfacePciE
	DiskInterfaceNvme
	DiskInterfaceUnknown
)

func DiskInterfaceFromString(v string) DiskInterface {
	switch strings.ToLower(v) {
	case "sata":
		return DiskInterfaceSata
	case "scsi":
		return DiskInterfaceScsi
	case "sas":
		return DiskInterfaceSas
	case "pcie":
		return DiskInterfacePciE
	case "nvme":
		return DiskInterfaceNvme
	default:
		return DiskInterfaceUnknown
	}
}

func (d DiskInterface) String() string {
	switch d {
	case DiskInterfaceSata:
		return "SATA"
	case DiskInterfaceScsi:
		return "SCSI"
	case DiskInterfaceSas:
		return "SAS"
	case DiskInterfacePciE:
		return "PCIe"
	case DiskInterfaceNvme:
		return "NVMe"
	default:
		return "Unknown"
	}
}
