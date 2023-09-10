$version: "2.0"

namespace awlsring.camp.common.ip

enum IpAddressVersion {
    V4
    V6
}

@pattern("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$")
@length(min: 7, max: 16)
string IpV4Address

structure IpAddressSummary {
    @required
    version: IpAddressVersion

    @required
    address: String
}

list IpAddressSummaries {
    member: IpAddressSummary
}
