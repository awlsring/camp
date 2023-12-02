$version: "2.0"

namespace awlsring.camp.local

@documentation("Defines a power capability for a machine.")
enum MachinePowerCapability {
    @documentation("The machine can be rebooted.")
    REBOOT = "Reboot"

    @documentation("The machine can be powered off.")
    POWER_OFF = "PowerOff"

    @documentation("The machine can be powered on via WakeOnLan.")
    WAKE_ON_LAN = "WakeOnLan"
}

@documentation("List of power capabilities")
list MachinePowerCapabilities {
    member: MachinePowerCapability
}
