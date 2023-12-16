$version: "2.0"

namespace awlsring.camp.common

@documentation("Information about the machine's network interfaces")
structure NetworkInterfaceSummary {
    @required
    name: String

    @required
    virtual: Boolean

    macAddress: String

    @required
    ipAddresses: IpAddressSummaries

    vendor: String

    duplex: String

    speed: String

    pciAddress: String
}

list NetworkInterfaceSummaries {
    member: NetworkInterfaceSummary
}
