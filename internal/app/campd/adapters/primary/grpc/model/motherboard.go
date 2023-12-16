package model

import (
	"github.com/awlsring/camp/internal/pkg/domain/motherboard"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	grpcmodel "github.com/awlsring/camp/pkg/grpc_model"
)

func MotherboardFromDomain(in *motherboard.Motherboard) *campd.MotherboardSummary {
	return &campd.MotherboardSummary{
		AssetTag: grpcmodel.NewStringValue(in.AssetTag),
		Product:  grpcmodel.NewStringValue(in.Product),
		Serial:   grpcmodel.NewStringValue(in.Serial),
		Vendor:   grpcmodel.NewStringValue(in.Vendor),
		Version:  grpcmodel.NewStringValue(in.Version),
	}
}
