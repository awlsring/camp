$version: "2.0"

metadata proto_options = [
    {
        "awlsring.camp.campd": {go_package: "\"github.com/awlsring/camp/campd\""}
        "awlsring.camp.common": {go_package: "\"github.com/awlsring/camp/campd\""}
        "aws.protocols": {go_package: "\"github.com/awlsring/camp/campd\""}
        smithytranslate: {go_package: "\"github.com/awlsring/camp/campd\""}
    }
]

namespace awlsring.camp.campd

use alloy.proto#protoEnabled
use awlsring.camp.common#Health
use awlsring.camp.common#UnauthorizedException
use aws.protocols#restJson1

@title("Campd")
@protoEnabled
@restJson1
@httpApiKeyAuth(name: "X-Api-Key", in: "header")
@paginated(inputToken: "nextToken", outputToken: "nextToken", pageSize: "pageSize")
service Campd {
    version: "2023-12-15"
    operations: [Health, Reboot, PowerOff, BootTime, Uptime, DescribeHost, DescribeBios, DescribeMotherboard, DescribeCpu, DescribeMemory, DescribeDisk, ListDisks, DescribeNetworkInterface, ListNetworkInterfaces, ListAddresses]
    errors: [UnauthorizedException]
}
