package model

import (
	"github.com/awlsring/camp/internal/pkg/domain/motherboard"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	grpcmodel "github.com/awlsring/camp/pkg/grpc_model"
)

func BiosFromDomain(in *motherboard.Bios) *campd.BiosSummary {
	return &campd.BiosSummary{
		Vendor:  grpcmodel.NewStringValue(in.Vendor),
		Version: grpcmodel.NewStringValue(in.Version),
		Date:    grpcmodel.NewStringValue(in.Date),
	}
}
