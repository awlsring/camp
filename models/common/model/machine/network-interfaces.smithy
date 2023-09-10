$version: "2.0"

namespace awlsring.camp.common.machine

use awlsring.camp.common.ip#IpAddressSummaries

@documentation("Information about the machine's network interfaces")
structure MachineNetworkInterfaceSummary {
    @required
    name: String

    @required
    addresses: IpAddressSummaries

    @required
    virtual: Boolean

    macAddress: String

    vendor: String

    mtu: Integer

    duplex: String

    speed: Integer
}

list MachineNetworkInterfaceSummaries {
    member: MachineNetworkInterfaceSummary
}
