$version: "2.0"

namespace awlsring.camp.local

use awlsring.camp.common.exceptions#UnauthorizedException
use awlsring.camp.common.operations#Health
use aws.protocols#restJson1

@title("Camp Local Controller")
@restJson1
@httpBearerAuth
@httpApiKeyAuth(name: "X-Api-Key", in: "header")
@auth([httpBearerAuth, httpApiKeyAuth])
@paginated(inputToken: "nextToken", outputToken: "nextToken", pageSize: "pageSize")
service CampLocal {
    version: "2022-10-20"
    resources: [Machine]
    operations: [Health, Heartbeat, Register, ReportStatusChange, ReportSystemChange]
    errors: [UnauthorizedException]
}
