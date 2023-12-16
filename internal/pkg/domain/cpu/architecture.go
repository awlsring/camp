package cpu

import (
	"strings"
)

type Architecture int64

const (
	Architecturex86 Architecture = iota
	ArchitectureArm64
	ArchitectureArmV7
	ArchitectureRiscV64
	ArchitectureUnknown
)

func ArchitectureFromString(v string) Architecture {
	switch strings.ToLower(v) {
	case "x86", "x86_64", "amd64":
		return Architecturex86
	case "arm7", "armv7", "armv7l":
		return ArchitectureArmV7
	case "arm8", "armv8", "armv8l", "aaarch64", "arm64":
		return ArchitectureArm64
	case "riscv64":
		return ArchitectureRiscV64
	default:
		return ArchitectureUnknown
	}
}

func (c Architecture) String() string {
	switch c {
	case Architecturex86:
		return "x86"
	case ArchitectureArmV7:
		return "armv7"
	case ArchitectureArm64:
		return "arm64"
	case ArchitectureRiscV64:
		return "riscv64"
	default:
		return "Unknown"
	}
}
