$version: "2.0"

namespace awlsring.camp.local

use smithy.framework#ValidationException

@documentation("Method called by agent machines to register with the control server. Internal use only.")
@http(method: "POST", uri: "/internal/register", code: 200)
operation Register {
    input: RegisterInput
    output: RegisterOutput
    errors: [
        ValidationException
    ]
}

@input
structure RegisterInput {
    @required
    summary: ReportedMachineSummary
}

@output
structure RegisterOutput {
    @required
    success: Boolean
}
