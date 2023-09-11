$version: "2.0"

namespace awlsring.camp.common.machine

enum DiskType {
    HDD = "HDD"
    SDD = "SSD"
    UNKNOWN = "Unknown"
}

enum DiskInterface {
    SATA = "SATA"
    SCSI = "SCSI"
    PCI_E = "PCIe"
    UNKNOWN = "Unknown"
}

@documentation("Information about the machine's disks")
structure MachineDiskSummary {
    @required
    device: String

    model: String

    vendor: String

    @required
    interface: DiskInterface

    serial: String

    @required
    type: DiskType

    sectorSize: Integer

    sizeRaw: Long

    @required
    size: Long
}

list MachineDiskSummaries {
    member: MachineDiskSummary
}
