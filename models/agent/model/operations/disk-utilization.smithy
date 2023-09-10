$version: "2.0"

namespace awlsring.camp.agent

use awlsring.camp.common.exceptions#ResourceNotFoundException
use smithy.framework#ValidationException

@documentation("Provides a summary of utilization of a disk")
@readonly
@http(method: "GET", uri: "/disk/{identifier}/utilization", code: 200)
operation GetDiskUtilization {
    input: GetDiskUtilizationInput
    output: GetDiskUtilizationOutput
    errors: [
        ValidationException
        ResourceNotFoundException
    ]
}

@input
structure GetDiskUtilizationInput {
    @httpLabel
    @required
    identifier: String
}

@output
structure GetDiskUtilizationOutput {
    @required
    summary: DiskUtilizationSummary
}

@documentation("Provides a summary of disk utilization.")
structure DiskUtilizationSummary {
    @required
    availableSpace: Long

    @required
    totalSpace: Long

    @required
    usedSpace: Long
}
