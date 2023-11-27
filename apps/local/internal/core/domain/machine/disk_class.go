package machine

import (
	"errors"
	"strings"
)

var (
	ErrInvalidDiskClass = errors.New("invalid disk class")
)

type DiskClass int64

const (
	DiskTypeHDD DiskClass = iota
	DiskTypeSSD
	DiskTypeUnknown
)

func DiskClassFromString(v string) DiskClass {
	switch strings.ToLower(v) {
	case "hdd":
		return DiskTypeHDD
	case "ssd":
		return DiskTypeSSD
	default:
		return DiskTypeUnknown
	}
}

func (d DiskClass) String() string {
	switch d {
	case DiskTypeHDD:
		return "HDD"
	case DiskTypeSSD:
		return "SSD"
	default:
		return "Unknown"
	}
}
