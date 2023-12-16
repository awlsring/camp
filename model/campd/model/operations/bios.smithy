$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common#BiosSummary

@documentation("Returns a summary of the machine's BIOS.")
@readonly
@http(method: "GET", uri: "/bios", code: 200)
operation DescribeBios {
    output := {
        @required
        bios: BiosSummary
    }
}
