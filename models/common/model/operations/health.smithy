$version: "2.0"

namespace awlsring.camp.common.operations

@readonly
@http(method: "GET", uri: "/health", code: 200)
operation Health {
    output: HealthOutput
}

@output
structure HealthOutput {
    @required
    success: Boolean
}
