package machine

type InternalIdentifier string

func InternalIdentifierFromString(v string) (InternalIdentifier, error) {
	return NonNullString[InternalIdentifier](v)
}

func (i InternalIdentifier) String() string {
	return string(i)
}
