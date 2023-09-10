$version: "2.0"

namespace awlsring.camp.agent

@documentation("Provides an overview of the machine's resources.")
@readonly
@http(method: "GET", uri: "/uptime", code: 200)
operation GetUptime {
    input: GetUptimeInput
    output: GetUptimeOutput
}

@input
structure GetUptimeInput {}

@output
structure GetUptimeOutput {
    @required
    summary: UptimeSummary
}

@documentation("Provides an overview of the machines uptime.")
structure UptimeSummary {
    @documentation("The machine's uptime in seconds.")
    @required
    upTime: Long

    @documentation("The machine's boot time.")
    @required
    bootTime: Timestamp
}
