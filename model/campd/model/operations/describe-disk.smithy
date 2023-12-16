$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common#DiskSummary
use awlsring.camp.common#ResourceNotFoundException

@documentation("Returns a summary a specified disk.")
@readonly
@http(method: "GET", uri: "/disk/{name}", code: 200)
operation DescribeDisk {
    input := {
        @required
        @httpLabel
        name: String
    }

    output := {
        @required
        disk: DiskSummary
    }

    errors: [
        ResourceNotFoundException
    ]
}
