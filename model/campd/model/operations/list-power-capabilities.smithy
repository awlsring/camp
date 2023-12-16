$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common.power#PowerCapabilities

@documentation("Returns a summary of the machine's BIOS.")
@readonly
@http(method: "GET", uri: "/power", code: 200)
operation ListPowerCapabilities {
    output := {
        @required
        capabilities: PowerCapabilities
    }
}
