$version: "2.0"

namespace awlsring.camp.local

@mixin
structure MachinePowerCapabilityCommonMixin {
    @documentation("If this cabapility is enabled")
    @required
    enabled: Boolean
}

@documentation("Defines the reboot capability for a machine.")
structure MachinePowerCapabilityRebootSummary with [MachinePowerCapabilityCommonMixin] {}

@documentation("Defines the power off capability for a machine.")
structure MachinePowerCapabilityPowerOffSummary with [MachinePowerCapabilityCommonMixin] {}

@documentation("Defines the WakeOnLan capability for a machine.")
structure MachinePowerCapabilityWakeOnLanSummary with [MachinePowerCapabilityCommonMixin] {
    @documentation("The MAC address of the machine. This is required if the capability is enabled.")
    macAddress: String
}

@documentation("Defines all power capabilities for a machine.")
structure MachinePowerCapabilitiesSummary {
    @required
    reboot: MachinePowerCapabilityRebootSummary

    @required
    powerOff: MachinePowerCapabilityPowerOffSummary

    @required
    wakeOnLan: MachinePowerCapabilityWakeOnLanSummary
}
