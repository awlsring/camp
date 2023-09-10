$version: "2.0"

namespace awlsring.camp.agent

use awlsring.camp.common.exceptions#ResourceNotFoundException
use smithy.framework#ValidationException

@documentation("Provides a summary of utilization of a volume")
@readonly
@http(method: "GET", uri: "/volume/{identifier}/utilization", code: 200)
operation GetVolumeUtilization {
    input: GetVolumeUtilizationInput
    output: GetVolumeUtilizationOutput
    errors: [
        ValidationException
        ResourceNotFoundException
    ]
}

@input
structure GetVolumeUtilizationInput {
    @httpLabel
    @required
    identifier: String
}

@output
structure GetVolumeUtilizationOutput {
    @required
    summary: VolumeUtilizationSummary
}

@documentation("Provides a summary of volume utilization.")
structure VolumeUtilizationSummary {
    @required
    availableSpace: Long

    @required
    totalSpace: Long

    @required
    usedSpace: Long
}
