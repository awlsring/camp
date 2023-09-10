$version: "2.0"

namespace awlsring.camp.common.machine

@pattern("^m-[a-zA-Z0-9\b]{32}$")
string MachineId

enum MachineStatus {
    RUNNING = "Running"
    STARTING = "Starting"
    STOPPING = "Stopping"
    STOPPED = "Stopped"
    UNKNOWN = "Unknown"
}

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
