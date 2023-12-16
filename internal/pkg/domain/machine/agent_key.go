package machine

import "github.com/awlsring/camp/internal/pkg/values"

type AgentKey string

func AgentKeyFromString(v string) (AgentKey, error) {
	return values.NonNullString[AgentKey](v)
}

func (m AgentKey) String() string {
	return string(m)
}
