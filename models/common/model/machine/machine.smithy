$version: "2.0"

namespace awlsring.camp.common.machine

@pattern("^m-[a-zA-Z0-9\b]{32}$")
string MachineId

@documentation("The power status of a machine")
enum MachineStatus {
    RUNNING = "Running"
    STARTING = "Starting"
    STOPPING = "Stopping"
    REBOOTING = "Rebooting"
    STOPPED = "Stopped"
    PENDING = "Pending"
    UNKNOWN = "Unknown"
}

@documentation("The class of a machine.")
enum MachineClass {
    BARE_METAL = "BareMetal"
    VIRTUAL_MACHINE = "VirtualMachine"
    HYPERVISOR = "Hypervisor"
    UNKNOWN = "Unknown"
}

structure MachineStatusSummary {
    @documentation("The machine's last reported status")
    @required
    status: MachineStatus

    @documentation("Timestamp of the last status check")
    @required
    lastChecked: Timestamp
}
