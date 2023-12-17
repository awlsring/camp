package reporter

import (
	"github.com/awlsring/camp/internal/app/campd/ports/reporting"
	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
)

type Service struct {
	id        machine.Identifier
	reporting reporting.Reporting
	hostSvc   service.Host
}
