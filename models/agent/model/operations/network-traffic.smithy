$version: "2.0"

namespace awlsring.camp.agent

use awlsring.camp.common.exceptions#ResourceNotFoundException
use smithy.framework#ValidationException

@documentation("Provides a summary of utilization of a disk")
@readonly
@http(method: "GET", uri: "/network/{identifier}/utilization", code: 200)
operation GetNetworkInterfaceUtilization {
    input: GetNetworkInterfaceUtilizationInput
    output: GetNetworkInterfaceUtilizationOutput
    errors: [
        ValidationException
        ResourceNotFoundException
    ]
}

@input
structure GetNetworkInterfaceUtilizationInput {
    @httpLabel
    @required
    identifier: String
}

@output
structure GetNetworkInterfaceUtilizationOutput {
    @required
    summary: NetworkInterfaceUtilizationSummary
}

@documentation("Provides a summary of disk utilization.")
structure NetworkInterfaceUtilizationSummary {
    @required
    bytesTraffic: NetworkInterfaceTrafficSummary

    @required
    packetTraffic: NetworkInterfaceTrafficSummary
}

structure NetworkInterfaceTrafficSummary {
    @required
    transmitted: Long

    @required
    recieved: Long
}
