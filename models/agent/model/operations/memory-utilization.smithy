$version: "2.0"

namespace awlsring.camp.agent

use smithy.framework#ValidationException

@documentation("Provides a summary of utilization of system memory.")
@readonly
@http(method: "GET", uri: "/memory/utilization", code: 200)
operation GetMemoryUtilization {
    input: GetMemoryUtilizationInput
    output: GetMemoryUtilizationOutput
    errors: [
        ValidationException
    ]
}

@input
structure GetMemoryUtilizationInput {}

@output
structure GetMemoryUtilizationOutput {
    @required
    summary: MemoryUtilizationSummary
}

@documentation("Provides a summary of volume utilization.")
structure MemoryUtilizationSummary {
    @required
    available: Long

    @required
    total: Long

    @required
    used: Long
}
