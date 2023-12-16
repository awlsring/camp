$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common#ResourceNotFoundException
use smithy.framework#ValidationException

@documentation("Method called by agent machines to report a to their system. Internal use only.")
@http(method: "POST", uri: "/internal/change/system", code: 200)
operation ReportSystemChange {
    input := {
        @required
        @required
        identifier: InternalMachineIdentifier

        @required
        summary: ReportedMachineSummary
    }

    output := {
        @required
        success: Boolean
    }

    errors: [
        ValidationException
        ResourceNotFoundException
    ]
}
