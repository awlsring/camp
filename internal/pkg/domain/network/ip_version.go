package network

import (
	"errors"
	"strings"
)

var (
	ErrInvalidIpAddressClass = errors.New("invalid ip address class")
)

type IpAddressVersion int64

const (
	IpAddressV4 IpAddressVersion = iota
	IpAddressV6
	IpAddressTypeUnknown
)

func DetermineIpAddressType(ip string) IpAddressVersion {
	if strings.Contains(ip, ":") {
		return IpAddressV6
	}
	if strings.Contains(ip, ".") {
		return IpAddressV4
	}
	return IpAddressTypeUnknown
}

func IpAddressTypeFromString(v string) IpAddressVersion {
	switch strings.ToLower(v) {
	case "ipv4", "v4":
		return IpAddressV4
	case "ipv6", "v6":
		return IpAddressV6
	default:
		return IpAddressTypeUnknown
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
