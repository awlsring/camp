$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common#CapabilityNotEnabledException
use awlsring.camp.common#InvalidPowerStateException
use awlsring.camp.common#ResourceNotFoundException
use awlsring.camp.common#StatusCode
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
        status: StatusCode
    }

    errors: [
        ValidationException
        ResourceNotFoundException
        InvalidPowerStateException
        CapabilityNotEnabledException
    ]
}
