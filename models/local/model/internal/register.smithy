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
    @documentation("The summary of the machine to register.")
    @required
    summary: ReportedMachineSummary

    @documentation("The endpoint to use for callbacks.")
    @required
    callbackEndpoint: String

    @documentation("The key to use for callbacks.")
    @required
    callbackKey: String
}

@output
structure RegisterOutput {
    @required
    success: Boolean
}
