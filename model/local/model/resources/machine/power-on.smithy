$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common#CapabilityNotEnabledException
use awlsring.camp.common#InvalidPowerStateException
use awlsring.camp.common#ResourceNotFoundException
use awlsring.camp.common#StatusCode
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
        status: StatusCode
    }

    errors: [
        ValidationException
        ResourceNotFoundException
        InvalidPowerStateException
        CapabilityNotEnabledException
    ]
}
