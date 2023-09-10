$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common.ip#IpAddressSummaries
use awlsring.camp.common.machine#MachineClass
use awlsring.camp.common.machine#MachineCpuSummary
use awlsring.camp.common.machine#MachineDiskSummaries
use awlsring.camp.common.machine#MachineMemorySummary
use awlsring.camp.common.machine#MachineNetworkInterfaceSummaries
use awlsring.camp.common.machine#MachineSystemSummary
use awlsring.camp.common.machine#MachineVolumeSummaries

@documentation("Summarized information about a machine")
structure ReportedMachineSummary {
    @documentation("The machine's internal identifier")
    @required
    internalIdentifier: String

    @documentation("The machine's class")
    class: MachineClass

    @documentation("Information about the machine system")
    @required
    system: MachineSystemSummary

    @documentation("CPU information")
    @required
    cpu: MachineCpuSummary

    @documentation("Memory information")
    @required
    memory: MachineMemorySummary

    @documentation("Disk information")
    @required
    disks: MachineDiskSummaries

    @documentation("Volume information")
    @required
    volumes: MachineVolumeSummaries

    @documentation("Network information")
    @required
    networkInterfaces: MachineNetworkInterfaceSummaries

    @documentation("IP Address information")
    @required
    addresses: IpAddressSummaries
}
