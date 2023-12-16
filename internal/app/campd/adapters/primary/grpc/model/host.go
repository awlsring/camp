package model

import (
	"github.com/awlsring/camp/internal/pkg/domain/host"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	grpcmodel "github.com/awlsring/camp/pkg/grpc_model"
)

func HostFromDomain(in *host.Host) *campd.HostSummary {
	os := OsFromDomain(in.OS)
	return &campd.HostSummary{
		Hostname: grpcmodel.NewStringValue(in.Hostname),
		HostId:   grpcmodel.NewStringValue(in.HostId),
		Os:       os,
	}
}

func OsFromDomain(in *host.OS) *campd.OsSummary {
	return &campd.OsSummary{
		Kernel:   grpcmodel.NewStringValue(in.Kernel),
		Name:     grpcmodel.NewStringValue(in.Name),
		Family:   grpcmodel.NewStringValue(in.Family),
		Platform: grpcmodel.NewStringValue(in.Platform),
		Version:  grpcmodel.NewStringValue(in.Version),
	}
}
