$version: "2.0"

namespace awlsring.camp.common.machine

@documentation("Information about the machine system")
structure MachineSystemSummary {
    family: String
    kernelVersion: String
    os: String
    osVersion: String
    osPretty: String
    hostname: String
}
