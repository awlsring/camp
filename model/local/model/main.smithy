$version: "2.0"

metadata proto_options = [
    {
        "awlsring.camp.local": {go_package: "\"github.com/awlsring/camp/local\""}
        "awlsring.camp.common": {go_package: "\"github.com/awlsring/camp/local\""}
        "aws.protocols": {go_package: "\"github.com/awlsring/camp/local\""}
        "smithy.framework": {go_package: "\"github.com/awlsring/camp/local\""}
        smithytranslate: {go_package: "\"github.com/awlsring/camp/local\""}
    }
]

namespace awlsring.camp.local

use alloy.proto#protoEnabled
use awlsring.camp.common#Health
use awlsring.camp.common#UnauthorizedException
use aws.protocols#restJson1

@title("CampLocal")
@protoEnabled
@restJson1
@httpApiKeyAuth(name: "X-Api-Key", in: "header")
@paginated(inputToken: "nextToken", outputToken: "nextToken", pageSize: "pageSize")
service CampLocal {
    version: "2023-12-15"
    resources: [Machine]
    operations: [Health, Heartbeat, Register, ReportStatusChange, ReportSystemChange]
    errors: [UnauthorizedException]
}
