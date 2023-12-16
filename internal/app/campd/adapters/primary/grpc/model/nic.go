package model

import (
	"github.com/awlsring/camp/internal/pkg/domain/network"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	grpcmodel "github.com/awlsring/camp/pkg/grpc_model"
)

func NicFromDomain(in *network.Nic) *campd.NetworkInterfaceSummary {
	address := []*campd.IpAddressSummary{}
	for _, a := range in.IpAddresses {
		address = append(address, IpAddressFromDomain(a))
	}
	mac := in.MacAddress.String()
	return &campd.NetworkInterfaceSummary{
		Name:        in.Name,
		Virtual:     in.Virtual,
		MacAddress:  grpcmodel.NewStringValue(&mac),
		Vendor:      grpcmodel.NewStringValue(in.Vendor),
		Duplex:      grpcmodel.NewStringValue(in.Duplex),
		Speed:       grpcmodel.NewStringValue(in.Speed),
		PCIAddress:  in.PCIAddress,
		IpAddresses: address,
	}
}
