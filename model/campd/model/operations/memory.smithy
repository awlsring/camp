$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common#MemorySummary

@documentation("Returns a summary of the machine's memory.")
@readonly
@http(method: "GET", uri: "/memory", code: 200)
operation DescribeMemory {
    output := {
        @required
        memory: MemorySummary
    }
}
