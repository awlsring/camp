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

@documentation("Defines a power capability for a machine.")
enum MachinePowerCapabilityType {
    @documentation("The machine can be rebooted.")
    REBOOT = "Reboot"

    @documentation("The machine can be powered off.")
    POWER_OFF = "PowerOff"

    @documentation("The machine can be powered on via WakeOnLan.")
    WAKE_ON_LAN = "WakeOnLan"
}

@documentation("List of power capabilities")
list MachinePowerCapabilitiesTypes {
    member: MachinePowerCapabilityType
}
