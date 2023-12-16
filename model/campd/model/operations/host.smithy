$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common#HostSummary

@documentation("Gets a description of host level details.")
@readonly
@http(method: "GET", uri: "/host", code: 200)
operation DescribeHost {
    output := {
        @required
        host: HostSummary
    }
}
