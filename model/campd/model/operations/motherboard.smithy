$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common#MotherboardSummary

@documentation("Returns a summary of the machine's motherboard.")
@readonly
@http(method: "GET", uri: "/mobo", code: 200)
operation DescribeMotherboard {
    output := {
        @required
        motherboard: MotherboardSummary
    }
}
