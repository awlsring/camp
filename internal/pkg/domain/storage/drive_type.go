package storage

import "strings"

type DriveType int

const (
	DriveTypeUnknown DriveType = iota
	DriveTypeHDD
	DriveTypeSSD
	DriveTypeFDD
	DriveTypeODD
	DriveTypeVirtual
)

func (d DriveType) String() string {
	switch d {
	case DriveTypeHDD:
		return "HDD"
	case DriveTypeSSD:
		return "SSD"
	case DriveTypeFDD:
		return "FDD"
	case DriveTypeODD:
		return "ODD"
	case DriveTypeVirtual:
		return "Virtual"
	default:
		return "Unknown"
	}
}

func DriveTypeFromString(v string) DriveType {
	switch strings.ToLower(v) {
	case "hdd":
		return DriveTypeHDD
	case "ssd":
		return DriveTypeSSD
	case "fdd":
		return DriveTypeFDD
	case "odd":
		return DriveTypeODD
	case "virtual", "virt":
		return DriveTypeVirtual
	default:
		return DriveTypeUnknown
	}
}
