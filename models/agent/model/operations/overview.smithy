$version: "2.0"

namespace awlsring.camp.agent

use awlsring.camp.common.ip#IpAddressSummaries
use awlsring.camp.common.machine#MachineCpuSummary
use awlsring.camp.common.machine#MachineDiskSummaries
use awlsring.camp.common.machine#MachineMemorySummary
use awlsring.camp.common.machine#MachineNetworkInterfaceSummaries
use awlsring.camp.common.machine#MachineSystemSummary
use awlsring.camp.common.machine#MachineVolumeSummaries

@documentation("Provides an overview of the machine's resources.")
@readonly
@http(method: "GET", uri: "/overview", code: 200)
operation GetOverview {
    input: GetOverviewInput
    output: GetOverviewOutput
}

@input
structure GetOverviewInput {}

@output
structure GetOverviewOutput {
    @required
    summary: OverviewSummary
}

structure OverviewSummary {
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
