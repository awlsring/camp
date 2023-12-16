$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common#IpAddressSummaries

@documentation("Returns summaries of available ip addresses on the machine.")
@readonly
@http(method: "GET", uri: "/address", code: 200)
operation ListAddresses {
    output := {
        @required
        addresses: IpAddressSummaries
    }
}
