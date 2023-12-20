package network

import "github.com/awlsring/camp/internal/pkg/values"

type IpAddress struct {
	Version IpAddressVersion
	Address Address
	Nic     *Nic
}

func NewIpAddress(version IpAddressVersion, address Address, nic *Nic) IpAddress {
	return IpAddress{
		Version: version,
		Address: address,
		Nic:     nic,
	}
}

type Address string

func (a Address) String() string {
	return string(a)
}

func AddressFromString(v string) (Address, error) {
	return values.NonNullString[Address](v)
}
