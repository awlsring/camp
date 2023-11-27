package machine

type DiskIdentifier string

func (d DiskIdentifier) String() string {
	return string(d)
}

func DiskIdentifierFromString(v string) (DiskIdentifier, error) {
	return NonNullString[DiskIdentifier](v)
}
