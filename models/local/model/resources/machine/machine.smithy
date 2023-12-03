$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common.ip#IpAddressSummaries
use awlsring.camp.common.machine#MachineClass
use awlsring.camp.common.machine#MachineCpuSummary
use awlsring.camp.common.machine#MachineDiskSummaries
use awlsring.camp.common.machine#MachineMemorySummary
use awlsring.camp.common.machine#MachineNetworkInterfaceSummaries
use awlsring.camp.common.machine#MachineStatusSummary
use awlsring.camp.common.machine#MachineSystemSummary
use awlsring.camp.common.machine#MachineVolumeSummaries
use awlsring.camp.common.tags#Tags

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

    @documentation("Information about the machine status")
    @required
    status: MachineStatusSummary

    @documentation("The machine's class")
    class: MachineClass

    @documentation("Tags attached to the machine")
    tags: Tags

    @documentation("The power state capabilities of a machine.")
    @required
    powerCapabilities: MachinePowerCapabilitiesSummary

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

@documentation("List of machine summaries")
list MachineSummaries {
    member: MachineSummary
}
