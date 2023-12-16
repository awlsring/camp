$version: "2.0"

namespace awlsring.camp.common

structure IpAddressSummary {
    @required
    version: IpAddressVersion

    @required
    address: String
}

list IpAddressSummaries {
    member: IpAddressSummary
}
