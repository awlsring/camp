$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common.exceptions#ResourceNotFoundException
use smithy.framework#ValidationException

@documentation(
    "Powers on a machine. Requires the machine has a capability that allows `PowerOn` (such as `WakeOnLan`) and is currently in a Stopped state."
)
@http(method: "POST", uri: "/machine/{identifier}/poweron", code: 200)
operation PowerOnMachine {
    input := {
        @httpLabel
        @required
        identifier: InternalMachineId
    }

    output := {
        @required
        success: Boolean
    }

    errors: [
        ValidationException
        ResourceNotFoundException
        InvalidPowerStateException
        CapabilityNotEnabledException
    ]
}
