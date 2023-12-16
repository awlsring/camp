$version: "2.0"

namespace awlsring.camp.campd

@documentation("Get the machines boot time.")
@readonly
@http(method: "GET", uri: "/boottime", code: 200)
operation BootTime {
    output := {
        @required
        bootTime: Timestamp
    }
}
