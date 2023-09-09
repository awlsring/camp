$version: "2.0"

namespace awlsring.camp.local

use aws.protocols#restJson1
@title("CampLocal")

@restJson1
@httpBearerAuth
@httpApiKeyAuth(name: "X-Api-Key", in: "header")
service CampLocal {
    version: "2022-10-20",
    operations: [Health]
}