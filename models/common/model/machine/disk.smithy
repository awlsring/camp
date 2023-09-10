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

    @required
    model: String

    @required
    vendor: String

    @required
    interface: DiskInterface

    @required
    serial: String

    @required
    type: DiskType

    @required
    sectorSize: Integer

    @required
    sizeRaw: Long
}

list MachineDiskSummaries {
    member: MachineDiskSummary
}
