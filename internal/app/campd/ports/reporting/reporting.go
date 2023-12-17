package reporting

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/core/domain/system"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/power"
)

type Reporting interface {
	Heartbeat(ctx context.Context, id machine.Identifier) error
	Register(ctx context.Context, id machine.Identifier, class machine.MachineClass, system *system.System, power *machine.PowerCapabilities, callback machine.MachineEndpoint) error
	ReportStatus(ctx context.Context, id machine.Identifier, status power.StatusCode) error
	ReportSystemInformation(ctx context.Context, id machine.Identifier, system *system.System) error
}
