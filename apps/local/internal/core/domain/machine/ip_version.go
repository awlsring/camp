package machine

import (
	"errors"
	"strings"
)

var (
	ErrInvalidIpAddressClass = errors.New("invalid id address class")
)

type IpAddressVersion int64

const (
	IpAddressV4 IpAddressVersion = iota
	IpAddressV6
	IpAddressTypeUnknown
)

func IpAddressTypeFromString(v string) (IpAddressVersion, error) {
	switch strings.ToLower(v) {
	case "ipv4", "v4":
		return IpAddressV4, nil
	case "ipv6", "v6":
		return IpAddressV6, nil
	default:
		return IpAddressTypeUnknown, ErrInvalidIpAddressClass
	}
}

func (i IpAddressVersion) String() string {
	switch i {
	case IpAddressV4:
		return "IPv4"
	case IpAddressV6:
		return "IPv6"
	default:
		return "Unknown"
	}
}
