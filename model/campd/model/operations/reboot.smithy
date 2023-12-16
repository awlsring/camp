$version: "2.0"

namespace awlsring.camp.campd

use awlsring.camp.common#CapabilityNotEnabledException

@documentation("Reboots the machine.")
@http(method: "POST", uri: "/reboot", code: 200)
operation Reboot {
    output := {
        success: Boolean
    }
    errors: [
        CapabilityNotEnabledException
    ]
}
