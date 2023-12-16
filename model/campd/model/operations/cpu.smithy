$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common#CpuSummary

@documentation("Returns a summary of the machine's CPUs.")
@readonly
@http(method: "GET", uri: "/cpu", code: 200)
operation DescribeCpu {
    output := {
        @required
        cpu: CpuSummary
    }
}
