$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common#CpuSummary
use awlsring.camp.common#DiskSummaries
use awlsring.camp.common#HostSummary
use awlsring.camp.common#IpAddressSummaries
use awlsring.camp.common#MachineClass
use awlsring.camp.common#MemorySummary
use awlsring.camp.common#NetworkInterfaceSummaries
use smithy.framework#ValidationException

@documentation("Method called by agent machines to register with the control server. Internal use only.")
@http(method: "POST", uri: "/internal/register", code: 200)
operation Register {
    input := {
        @documentation("The machine's internal identifier")
        @required
        internalIdentifier: InternalMachineIdentifier

        @documentation("The machine's class")
        @required
        class: MachineClass

        @documentation("The summary of the machine to register.")
        @required
        systemSummary: ReportedMachineSummary

        @documentation("Description of power capabilities of the machine.")
        powerCapabilities: ReportedPowerCapabilitiesSummary

        @documentation("The endpoint to use for callbacks.")
        @required
        callbackEndpoint: String
    }

    output := {
        @documentation("The key to use for callbacks.")
        @required
        accessKey: String
    }

    errors: [
        ValidationException
    ]
}

@documentation("Defines the reported power capabilities for a machine.")
structure ReportedPowerCapabilitiesSummary {
    reboot: MachinePowerCapabilityRebootSummary
    powerOff: MachinePowerCapabilityPowerOffSummary
    wakeOnLan: MachinePowerCapabilityWakeOnLanSummary
}

@documentation("Summarized information about a machine")
structure ReportedMachineSummary {
    @documentation("Information about the machine system")
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
