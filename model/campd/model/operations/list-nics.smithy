$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common#NetworkInterfaceSummaries

@documentation("Returns summaries for all network interfaces on the machine.")
@readonly
@http(method: "GET", uri: "/nic", code: 200)
operation ListNetworkInterfaces {
    output := {
        @required
        nics: NetworkInterfaceSummaries
    }
}
