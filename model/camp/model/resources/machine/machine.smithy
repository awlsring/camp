$version: "2.0"

namespace awlsring.camp.api

use awlsring.camp.common#IpAddressSummaries
use awlsring.camp.common#CpuSummary
use awlsring.camp.common#DiskSummaries
use awlsring.camp.common#HostSummary
use awlsring.camp.common#MachineClass
use awlsring.camp.common#StatusSummary
use awlsring.camp.common#MemorySummary
use awlsring.camp.common#NetworkInterfaceSummaries
use awlsring.camp.common#Tags

resource Machine {
    identifiers: {identifier: MachineIdentifier}
}

@documentation("The machine's identifier.")
@length(min: 1, max: 128)
string MachineIdentifier

@documentation("Summarized information about a machine")
structure MachineSummary {
    @documentation("The machine identifier")
    @required
    identifier: MachineIdentifier

    @documentation("Information about the machine status")
    @required
    status: MachineStatusSummary

    @documentation("The machine's class")
    class: MachineClass

    @documentation("Tags attached to the machine")
    tags: Tags

    // @documentation("The power state capabilities of a machine.")
    // @required
    // powerCapabilities: MachinePowerCapabilitiesSummary

    @documentation("Information on the machines site")
    @required
    site: MachineSiteSummary

    @documentation("Information about the host.")
    @required
    host: HostSummary

    @documentation("CPU information")
    @required
    cpu: CpuSummary

    @documentation("Memory information")
    @required
    memory: MemorySummary

    @required
    storage: StorageSummary

    @documentation("Network information")
    @required
    network: NetworkSummary
}

structure MachineSiteSummary {
    @documentation("The site identifier")
    @required
    identifier: SiteIdentifier

    @documentation("The site name")
    @required
    name: SiteName
}


@documentation("List of machine summaries")
list MachineSummaries {
    member: MachineSummary
}
