$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common#ResourceNotFoundException
use smithy.framework#ValidationException

@documentation("Describe a particular machine.")
@readonly
@http(method: "GET", uri: "/machine/{identifier}", code: 200)
operation DescribeMachine {
    input := {
        @httpLabel
        @required
        identifier: InternalMachineIdentifier
    }

    output := {
        @required
        summary: MachineSummary
    }

    errors: [
        ResourceNotFoundException
        ValidationException
    ]
}
