$version: "2.0"

namespace awlsring.camp.local

use smithy.framework#ValidationException

@documentation("List all machines that control server is aware of.")
@readonly
@http(method: "GET", uri: "/machine", code: 200)
operation ListMachines {
    output := {
        @required
        summaries: MachineSummaries
    }

    errors: [
        ValidationException
    ]
}
