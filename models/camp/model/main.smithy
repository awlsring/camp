$version: "2.0"

namespace awlsring.camp.api

use awlsring.camp.common.exceptions#UnauthorizedException
use awlsring.camp.common.operations#Health
use aws.protocols#restJson1

@title("Camp API Controller")
@restJson1
@httpBearerAuth
@httpApiKeyAuth(name: "X-Api-Key", in: "header")
@auth([httpBearerAuth, httpApiKeyAuth])
@paginated(inputToken: "nextToken", outputToken: "nextToken", pageSize: "pageSize")
service CampLocal {
    version: "2023-09-11"
    operations: [Health]
    errors: [UnauthorizedException]
}
