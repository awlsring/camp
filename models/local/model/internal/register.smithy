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
use smithy.framework#ValidationException

@documentation("Method called by agent machines to register with the control server. Internal use only.")
@http(method: "POST", uri: "/internal/register", code: 200)
operation Register {
    input := {
        @documentation("The machine's internal identifier")
        @required
        internalIdentifier: InternalMachineIdentifier

        @documentation("The machine's class")
        class: MachineClass

        @documentation("The summary of the machine to register.")
        @required
        systemSummary: ReportedMachineSummary

        @documentation("Description of power capabilities of the machine.")
        powerCapabilities: ReportedPowerCapabilitiesSummary

        @documentation("The endpoint to use for callbacks.")
        @required
        callbackEndpoint: String

        @documentation("The key to use for callbacks.")
        @required
        callbackKey: String
    }

    output := {
        @required
        success: Boolean
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
