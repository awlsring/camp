$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common.exceptions#ResourceNotFoundException
use awlsring.camp.common.machine#MachineId
use smithy.framework#ValidationException

@documentation("Describe a particular machine.")
@readonly
@http(method: "GET", uri: "/machine/{identifier}", code: 200)
operation DescribeMachine {
    input: DescribeMachineInput
    output: DescribeMachineOutput
    errors: [
        ResourceNotFoundException
        ValidationException
    ]
}

@input
structure DescribeMachineInput {
    @httpLabel
    @required
    identifier: MachineId
}

@output
structure DescribeMachineOutput {
    @required
    summary: MachineSummary
}
