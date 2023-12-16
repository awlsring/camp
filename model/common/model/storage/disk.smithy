$version: "2.0"

namespace awlsring.camp.common

enum DiskType {
    HDD = "HDD"
    SSD = "SSD"
    FDD = "FDD"
    ODD = "ODD"
    VIRTUAL = "Virtual"
    UNKNOWN = "Unknown"
}

enum StorageController {
    IDE = "IDE"
    SCSI = "SCSI"
    NVME = "NVMe"
    VIRTIO = "Virtio"
    MMC = "MMC"
    LOOP = "Loop"
    UNKNOWN = "Unknown"
}

string DiskIdentifier

@documentation("Information about the machine's disks")
structure DiskSummary {
    @required
    name: DiskIdentifier

    @required
    size: Long

    @required
    type: DiskType

    @required
    storageController: StorageController

    @required
    removable: Boolean

    partitions: PartitionIdentifiers

    serial: String

    model: String

    vendor: String

    wwn: String
}

list DiskSummaries {
    member: DiskSummary
}
