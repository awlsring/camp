$version: "2.0"

namespace awlsring.camp.agent

use awlsring.camp.common.exceptions#CapabilityNotEnabledException

@documentation("Powers off the machine.")
@http(method: "POST", uri: "/poweroff", code: 200)
operation PowerOff {
    input := {
        @documentation("Forces power off by ignorning graceful shutdown.")
        force: Boolean
    }
    errors: [
        CapabilityNotEnabledException,
    ]
}