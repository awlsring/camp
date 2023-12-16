$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common#ResourceNotFoundException
use awlsring.camp.common#StatusCode
use smithy.framework#ValidationException

@documentation("Method called by agent machines to report a change in their status. Internal use only.")
@http(method: "POST", uri: "/internal/change/status", code: 200)
operation ReportStatusChange {
    input := {
        @required
        internalIdentifier: String

        @required
        status: StatusCode
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
