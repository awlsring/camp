$version: "2.0"

namespace awlsring.camp.common.power

@documentation("Defines a power capability for a machine.")
enum PowerCapability {
    @documentation("The machine can be rebooted.")
    REBOOT = "Reboot"

    @documentation("The machine can be powered off.")
    POWER_OFF = "PowerOff"

    @documentation("The machine can be powered on via WakeOnLan.")
    POWER_ON = "PowerOn"
}

@documentation("List of power capabilities")
list PowerCapabilities {
    member: PowerCapability
}
