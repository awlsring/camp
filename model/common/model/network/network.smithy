$version: "2.0"

namespace awlsring.camp.common

structure NetworkSummary {
    @required
    ipAddresses: IpAddressSummaries

    networkInterfaces: NetworkInterfaceSummaries
}
