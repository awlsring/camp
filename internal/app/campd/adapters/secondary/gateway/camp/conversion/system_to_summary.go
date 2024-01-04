package conversion

import (
	"github.com/awlsring/camp/internal/app/campd/core/domain/system"
	local "github.com/awlsring/camp/pkg/gen/local_grpc"
)

func SystemFromDomain(in *system.System) *local.ReportedMachineSummary {
	return &local.ReportedMachineSummary{}
}
