$version: "2.0"

namespace awlsring.camp.common

string PartitionIdentifier

structure PartitionSummary {
    @required
    name: PartitionIdentifier

    @required
    size: Long

    @required
    disk: DiskIdentifier

    label: String

    readonly: Boolean

    type: String

    fileSystem: String

    uuid: String

    mountPoint: String
}

list PartitionSummaries {
    member: PartitionSummary
}

list PartitionIdentifiers {
    member: PartitionIdentifier
}
