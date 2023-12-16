$version: "2.0"

namespace awlsring.camp.campd

@documentation("Get the machines uptime.")
@readonly
@http(method: "GET", uri: "/uptime", code: 200)
operation Uptime {
    output := {
        @required
        uptime: Timestamp
    }
}
