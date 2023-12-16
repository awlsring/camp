$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common#NetworkInterfaceSummary
use awlsring.camp.common#ResourceNotFoundException

@documentation("Returns a summary of a specified network interface.")
@readonly
@http(method: "GET", uri: "/nic/{name}", code: 200)
operation DescribeNetworkInterface {
    input := {
        @required
        @httpLabel
        name: String
    }

    output := {
        @required
        nic: NetworkInterfaceSummary
    }

    errors: [
        ResourceNotFoundException
    ]
}
