package machine

type NetworkInterface struct {
	Name        NetworkInterfaceIdentifier
	IpAddresses []*IpAddress
	Virtual     bool
	MacAddress  *MacAddress
	Vendor      *string
	Mtu         *int
	Speed       *int
	Duplex      *string
}

func NewNetworkInterface(name NetworkInterfaceIdentifier, ipAddresses []*IpAddress, virtual bool, macAddress *MacAddress, vendor *string, mtu *int, speed *int, duplex *string) *NetworkInterface {
	return &NetworkInterface{
		Name:        name,
		IpAddresses: ipAddresses,
		Virtual:     virtual,
		MacAddress:  macAddress,
		Vendor:      vendor,
		Mtu:         mtu,
		Speed:       speed,
		Duplex:      duplex,
	}
}
