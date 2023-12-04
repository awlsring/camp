$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common.exceptions#CapabilityNotEnabledException
use awlsring.camp.common.exceptions#ResourceNotFoundException
use awlsring.camp.common.machine#MachineStatus
use smithy.framework#ValidationException

@documentation(
    "Powers off a machine. Requires the machine has the capability `PowerOff` and is currently in a running state."
)
@http(method: "POST", uri: "/machine/{identifier}/poweroff", code: 200)
operation PowerOffMachine {
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
