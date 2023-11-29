package tag

type ResourceType int

const (
	ResourceTypeMachine ResourceType = iota
	ResourceTypeUnknown
)

func ResourceTypeFromString(s string) ResourceType {
	switch s {
	case "Machine":
		return ResourceTypeMachine
	default:
		return ResourceTypeUnknown
	}
}

func (r ResourceType) String() string {
	switch r {
	case ResourceTypeMachine:
		return "Machine"
	default:
		return "Unknown"
	}
}
