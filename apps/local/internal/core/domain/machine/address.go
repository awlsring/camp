package machine

type IpAddress struct {
	Version IpAddressVersion
	Address Address
}

func NewIpAddress(version IpAddressVersion, address Address) IpAddress {
	return IpAddress{
		Version: version,
		Address: address,
	}
}

type Address string

func (a Address) String() string {
	return string(a)
}

func AddressFromString(v string) (Address, error) {
	return NonNullString[Address](v)
}
