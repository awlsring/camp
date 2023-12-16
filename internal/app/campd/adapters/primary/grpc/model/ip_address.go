package model

import (
	"github.com/awlsring/camp/internal/pkg/domain/network"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
)

func VersionFromDomain(in network.IpAddressVersion) campd.IpAddressVersion {
	switch in {
	case network.IpAddressV4:
		return campd.IpAddressVersion_V4
	case network.IpAddressV6:
		return campd.IpAddressVersion_V6
	default:
		return campd.IpAddressVersion_IPADDRESSVERSION_UNKNOWN
	}
}

func IpAddressFromDomain(in *network.IpAddress) *campd.IpAddressSummary {
	v := VersionFromDomain(in.Version)
	return &campd.IpAddressSummary{
		Version: v,
		Address: in.Address.String(),
	}
}
