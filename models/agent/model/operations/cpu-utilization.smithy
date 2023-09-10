$version: "2.0"

namespace awlsring.camp.agent

use smithy.framework#ValidationException

@documentation("Provides a summary of cpu utilization.")
@readonly
@http(method: "GET", uri: "/cpu/utilization", code: 200)
operation GetCpuUtilization {
    input: GetCpuUtilizationInput
    output: GetCpuUtilizationOutput
    errors: [
        ValidationException
    ]
}

@input
structure GetCpuUtilizationInput {}

@output
structure GetCpuUtilizationOutput {
    @required
    summary: CpuUtilizationSummary
}

@documentation("Provides a summary of disk utilization.")
structure CpuUtilizationSummary {
    @required
    coreUtilization: CoreUtilizationMap
}

map CoreUtilizationMap {
    key: String
    value: CoreUtilizationSummary
}

structure CoreUtilizationSummary {
    @required
    name: String

    @required
    usage: Float

    @required
    frequency: Float
}
