$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common#CpuSummary
use awlsring.camp.common#DiskSummaries
use awlsring.camp.common#HostSummary
use awlsring.camp.common#IpAddressSummaries
use awlsring.camp.common#MachineClass
use awlsring.camp.common#MemorySummary
use awlsring.camp.common#NetworkInterfaceSummaries
use awlsring.camp.common#StatusSummary
use awlsring.camp.common#Tags

resource Machine {
    identifiers: {identifier: InternalMachineIdentifier}
    read: DescribeMachine
    list: ListMachines
    operations: [SendMachineWakeOnLan, PowerOffMachine, RebootMachine, AddTagsToMachine, RemoveTagFromMachine]
}

@documentation("The machine's internal identifier.")
@length(min: 1, max: 128)
string InternalMachineIdentifier

@documentation("Summarized information about a machine")
structure MachineSummary {
    @documentation("The machine identifier")
    @required
    identifier: InternalMachineIdentifier

    @documentation("The time the machine was registered")
    @required
    registeredAt: Timestamp

    @documentation("The time the machine was last updated")
    @required
    updatedAt: Timestamp

    @documentation("Information about the machine status")
    @required
    status: StatusSummary

    @documentation("The machine's class")
    class: MachineClass

    @documentation("Tags attached to the machine")
    tags: Tags

    @documentation("The power state capabilities of a machine.")
    @required
    powerCapabilities: MachinePowerCapabilitiesSummary

    @documentation("Information about the host.")
    @required
    host: HostSummary

    @documentation("CPU information")
    @required
    cpu: CpuSummary

    @documentation("Memory information")
    @required
    memory: MemorySummary

    @documentation("Disk information")
    @required
    disks: DiskSummaries

    @documentation("Network information")
    @required
    networkInterfaces: NetworkInterfaceSummaries

    @documentation("IP Address information")
    @required
    addresses: IpAddressSummaries
}

@documentation("List of machine summaries")
list MachineSummaries {
    member: MachineSummary
}
