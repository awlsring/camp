$version: "2.0"

namespace awlsring.camp.common.machine

@documentation("Information about the machine system")
structure MachineSystemSummary {
    @required
    family: String

    @required
    kernelVersion: String

    @required
    os: String

    @required
    osVersion: String

    @required
    osPretty: String

    @required
    hostname: String
}
