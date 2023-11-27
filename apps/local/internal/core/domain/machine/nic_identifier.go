package machine

type NetworkInterfaceIdentifier string

func (n NetworkInterfaceIdentifier) String() string {
	return string(n)
}

func NetworkInterfaceIdentifierFromString(v string) (NetworkInterfaceIdentifier, error) {
	return NonNullString[NetworkInterfaceIdentifier](v)
}
