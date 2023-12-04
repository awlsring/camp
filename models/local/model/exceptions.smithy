$version: "2.0"

namespace awlsring.camp.local

@error("client")
@httpError(400)
structure InvalidPowerStateException {
    @required
    message: String
}

@error("client")
@httpError(400)
structure InvalidPowerStateException {
    @required
    message: String
}

@error("client")
@httpError(404)
structure ResourceNotFoundException {
    @required
    message: String
}
