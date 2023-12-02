$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common.exceptions#ResourceNotFoundException
use smithy.framework#ValidationException

@documentation(
    "Reboots a machine. Requires the machine has a capability that allows `Reboot` and is currently in a Running state."
)
@http(method: "POST", uri: "/machine/{identifier}/reboot", code: 200)
operation RebootMachine {
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
