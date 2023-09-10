$version: "2.0"

namespace awlsring.camp.local

use smithy.framework#ValidationException

@documentation("Method called by agent machines to report a to their system. Internal use only.")
@http(method: "POST", uri: "/internal/change/system", code: 200)
operation ReportSystemChange {
    input: ReportSystemChangeInput
    output: ReportSystemChangeOutput
    errors: [
        ValidationException
    ]
}

@input
structure ReportSystemChangeInput {
    @required
    summary: ReportedMachineSummary
}

@output
structure ReportSystemChangeOutput {
    @required
    success: Boolean
}