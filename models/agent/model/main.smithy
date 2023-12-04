$version: "2.0"

namespace awlsring.camp.agent

use awlsring.camp.common.exceptions#UnauthorizedException
use awlsring.camp.common.operations#Health
use aws.protocols#restJson1

@title("Camp Machine Agent")
@restJson1
@httpApiKeyAuth(name: "X-Api-Key", in: "header")
@paginated(inputToken: "nextToken", outputToken: "nextToken", pageSize: "pageSize")
service CampAgent {
    version: "2023-06-07"
    operations: [Health, GetOverview, GetUptime, GetCpuUtilization, GetMemoryUtilization, GetDiskUtilization, GetNetworkInterfaceUtilization, GetVolumeUtilization, Reboot, PowerOff]
    errors: [UnauthorizedException]
}
