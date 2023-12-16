package machine

import "github.com/awlsring/camp/internal/pkg/values"

type Identifier string

func IdentifierFromString(identifier string) (Identifier, error) {
	return values.NonNullString[Identifier](identifier)
}

func (m Identifier) String() string {
	return string(m)
}
