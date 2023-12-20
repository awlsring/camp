package values

import "strings"

func ParseOptional(value string) *string {
	if value == "" {
		return nil
	}
	if strings.Contains(strings.ToLower(value), "unknown") {
		return nil
	}
	return &value
}
