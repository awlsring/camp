$version: "2.0"

namespace awlsring.camp.local

use aws.protocols#restJson1
use awlsring.camp.common.exceptions#UnauthorizedException

@title("CampLocal")

@restJson1
@httpBearerAuth
@httpApiKeyAuth(name: "X-Api-Key", in: "header")
@auth([httpBearerAuth, httpApiKeyAuth])
service CampLocal {
    version: "2022-10-20",
    operations: [ Health ]
    errors: [ UnauthorizedException ]
}