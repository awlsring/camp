$version: "2.0"

namespace awlsring.camp.common

string VolumeIdentifier

structure VolumeSummary {
    @required
    name: VolumeIdentifier

    mountPoint: String

    size: Long

    fileSystem: String
}

list VolumeSummaries {
    member: VolumeSummary
}
