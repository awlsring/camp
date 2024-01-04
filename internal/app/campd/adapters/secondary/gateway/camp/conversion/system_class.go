package conversion

import (
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	local "github.com/awlsring/camp/pkg/gen/local_grpc"
)

func ClassFromDomain(c machine.MachineClass) local.MachineClass {
	switch c {
	case machine.MachineClassBareMetal:
		return local.MachineClass_BARE_METAL
	case machine.MachineClassVirtual:
		return local.MachineClass_VIRTUAL_MACHINE
	case machine.MachineClassHypervisor:
		return local.MachineClass_HYPERVISOR
	default:
		return local.MachineClass_MACHINECLASS_UNKNOWN
	}
}
