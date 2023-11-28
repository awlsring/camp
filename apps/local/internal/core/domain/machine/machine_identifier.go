package machine

type Identifier string

func IdentifierFromString(identifier string) (Identifier, error) {
	return NonNullString[Identifier](identifier)
}

func (m Identifier) String() string {
	return string(m)
}
