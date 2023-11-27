package machine

import (
	"errors"
	"strings"

	camperror "github.com/awlsring/camp/internal/pkg/errors"
)

var (
	ErrInvalidClass = errors.New("invalid machine class")
)

type MachineClass int64

const (
	MachineClassBareMetal MachineClass = iota
	MachineClassVirtual
	MachineClassHypervisor
	MachineClassUnknown
)

func MachineClassFromString(v string) (MachineClass, error) {
	switch strings.ToLower(v) {
	case "baremetal":
		return MachineClassBareMetal, nil
	case "virtualmachine":
		return MachineClassVirtual, nil
	case "hypervisor":
		return MachineClassHypervisor, nil
	default:
		return MachineClassUnknown, camperror.New(camperror.ErrValidation, ErrInvalidClass)
	}
}

func (m MachineClass) String() string {
	switch m {
	case MachineClassBareMetal:
		return "BareMetal"
	case MachineClassVirtual:
		return "VirtualMachine"
	case MachineClassHypervisor:
		return "Hypervisor"
	default:
		return "Unknown"
	}
}
