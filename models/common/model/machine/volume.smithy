$version: "2.0"

namespace awlsring.camp.common.machine

@documentation("Information about a machine volume")
structure MachineVolumeSummary {
    @required
    name: String

    @required
    mountPoint: String

    @required
    totalSpace: Long

    @required
    fileSystem: String
}

list MachineVolumeSummaries {
    member: MachineVolumeSummary
}
