$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common.machine#MachineStatus
use smithy.framework#ValidationException

@documentation("Method called by agent machines to report a change in their status. Internal use only.")
@http(method: "POST", uri: "/internal/change/status", code: 200)
operation ReportStatusChange {
    input: ReportStatusChangeInput
    output: ReportStatusChangeOutput
    errors: [
        ValidationException
    ]
}

@input
structure ReportStatusChangeInput {
    @required
    internalIdentifier: String

    @required
    summary: MachineStatus
}

@output
structure ReportStatusChangeOutput {
    @required
    success: Boolean
}