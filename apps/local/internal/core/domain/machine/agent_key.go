package machine

type AgentKey string

func AgentKeyFromString(v string) (AgentKey, error) {
	return NonNullString[AgentKey](v)
}

func (m AgentKey) String() string {
	return string(m)
}
