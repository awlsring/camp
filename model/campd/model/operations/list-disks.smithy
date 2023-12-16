$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common#DiskSummaries

@documentation("Returns summaries for all disks on the machine.")
@readonly
@http(method: "GET", uri: "/disk", code: 200)
operation ListDisks {
    output := {
        @required
        disks: DiskSummaries
    }
}
