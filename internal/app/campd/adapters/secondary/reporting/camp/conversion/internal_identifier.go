package conversion

import (
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	local "github.com/awlsring/camp/pkg/gen/local_grpc"
)

func InternalIdentifierFromDomain(id machine.Identifier) *local.InternalMachineIdentifier {
	return &local.InternalMachineIdentifier{
		Value: id.String(),
	}
}
