package camp_reporting

import (
	"github.com/awlsring/camp/internal/app/campd/ports/reporting"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	local "github.com/awlsring/camp/pkg/gen/local_grpc"
)

type CampLocalReporting struct {
	client local.CampLocalClient
}

func New(client local.CampLocalClient) reporting.Reporting {
	return &CampLocalReporting{client: client}
}

func identifierFromDomain(id machine.Identifier) string {
	return id.String()
}
