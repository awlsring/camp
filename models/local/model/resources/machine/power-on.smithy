$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common.exceptions#CapabilityNotEnabledException
use awlsring.camp.common.exceptions#ResourceNotFoundException
use awlsring.camp.common.machine#MachineStatus
use smithy.framework#ValidationException

@documentation(
    "Powers on a machine. Requires the machine has a capability of `WakeOnLan` and is currently in a Stopped state."
)
@http(method: "POST", uri: "/machine/{identifier}/wol", code: 200)
operation SendMachineWakeOnLan {
    input := {
        @httpLabel
        @required
        identifier: InternalMachineIdentifier
    }

    output := {
        @required
        status: MachineStatus
    }

    errors: [
        ValidationException
        ResourceNotFoundException
        InvalidPowerStateException
        CapabilityNotEnabledException
    ]
}
