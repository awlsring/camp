package machine

import "github.com/awlsring/camp/internal/pkg/values"

type MachineEndpoint string

func MachineEndpointFromString(v string) (MachineEndpoint, error) {
	return values.NonNullString[MachineEndpoint](v)
}

func (m MachineEndpoint) String() string {
	return string(m)
}
