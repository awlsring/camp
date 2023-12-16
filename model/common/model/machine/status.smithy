$version: "2.0"

namespace awlsring.camp.common

@documentation("The power status of a machine")
enum StatusCode {
    RUNNING = "Running"
    STARTING = "Starting"
    STOPPING = "Stopping"
    REBOOTING = "Rebooting"
    STOPPED = "Stopped"
    PENDING = "Pending"
    UNKNOWN = "Unknown"
}

structure StatusSummary {
    @documentation("The machine's last reported status")
    @required
    status: StatusCode

    @documentation("Timestamp of the last status update")
    @required
    lastUpdated: Timestamp
}
