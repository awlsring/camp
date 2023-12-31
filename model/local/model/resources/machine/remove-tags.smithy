$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common#ResourceNotFoundException
use awlsring.camp.common#TagKey
use smithy.framework#ValidationException

@idempotent
@documentation("Deletes a tag from a machine")
@http(method: "DELETE", uri: "/machine/{identifier}/tag/{tagKey}", code: 200)
operation RemoveTagFromMachine {
    input := {
        @httpLabel
        @required
        identifier: InternalMachineIdentifier

        @httpLabel
        @required
        tagKey: TagKey
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
