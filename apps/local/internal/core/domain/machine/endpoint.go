package machine

type MachineEndpoint string

func MachineEndpointFromString(v string) (MachineEndpoint, error) {
	return NonNullString[MachineEndpoint](v)
}

func (m MachineEndpoint) String() string {
	return string(m)
}
