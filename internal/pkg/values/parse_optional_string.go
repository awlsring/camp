package values

func ParseOptional(value string) *string {
	if value == "" {
		return nil
	}
	if value == "unknown" {
		return nil
	}
	return &value
}
