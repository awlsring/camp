$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common#ResourceNotFoundException
use awlsring.camp.common#Tags
use smithy.framework#ValidationException

@documentation("Adds tags to a machine. Supports up to 50 in one call.")
@http(method: "POST", uri: "/machine/{identifier}/tag", code: 200)
operation AddTagsToMachine {
    input := {
        @httpLabel
        @required
        identifier: InternalMachineIdentifier

        @required
        @length(min: 1, max: 50)
        tags: Tags
    }

    output := {
        @required
        sucess: Boolean
    }

    errors: [
        ValidationException
        ResourceNotFoundException
    ]
}
