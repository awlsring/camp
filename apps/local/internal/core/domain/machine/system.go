package machine

type System struct {
	Os       Os
	Hostname *string
}

type Os struct {
	Kernel     *string
	Name       *string
	PrettyName *string
	Version    *string
	Family     *string
}

func NewSystem(hostname, kernel, name, pretty, version, family *string) System {
	return System{
		Os: Os{
			Kernel:     kernel,
			Name:       name,
			PrettyName: pretty,
			Version:    version,
			Family:     family,
		},
		Hostname: hostname,
	}
}
