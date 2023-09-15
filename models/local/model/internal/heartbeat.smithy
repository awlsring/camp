$version: "2.0"

namespace awlsring.camp.local

use smithy.framework#ValidationException
use awlsring.camp.common.exceptions#ResourceNotFoundException

@documentation("Method called by agent machines to indicate they are still running. Only used internally.")
@http(method: "POST", uri: "/internal/heartbeat", code: 200)
operation Heartbeat {
    input: HeartbeatInput
    output: HeartbeatOutput
    errors: [
        ValidationException
        ResourceNotFoundException
    ]
}

@input
structure HeartbeatInput {
    @required
    internalIdentifier: String
}

@output
structure HeartbeatOutput {
    @required
    success: Boolean
}
