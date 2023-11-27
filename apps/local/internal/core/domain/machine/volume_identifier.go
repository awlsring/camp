package machine

type VolumeIdentifier string

func (v VolumeIdentifier) String() string {
	return string(v)
}

func VolumeIdentifierFromString(v string) (VolumeIdentifier, error) {
	return NonNullString[VolumeIdentifier](v)
}
