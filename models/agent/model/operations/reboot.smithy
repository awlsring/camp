$version: "2.0"

namespace awlsring.camp.agent

use awlsring.camp.common.exceptions#CapabilityNotEnabledException

@documentation("Reboots the machine.")
@http(method: "POST", uri: "/reboot", code: 200)
operation Reboot {
    errors: [
        CapabilityNotEnabledException,
    ]
}