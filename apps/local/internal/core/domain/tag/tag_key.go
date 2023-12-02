package tag

import (
	"fmt"
	"regexp"
)

const (
	TagKeyMinLength = 1
	TagKeyMaxLength = 50
	TagKeyPattern   = "^[a-zA-Z0-9_]+( [a-zA-Z0-9_]+){0,127}$"
)

var (
	ErrTagKeyInvalidLength   = fmt.Errorf("tag key must be between 1 and 50 characters")
	ErrTagKeyPatternMismatch = fmt.Errorf("tag key must match pattern %s", TagKeyPattern)
)

type TagKey string

func (t TagKey) String() string {
	return string(t)
}

func TagKeyFromString(s string) (TagKey, error) {
	if len(s) < TagKeyMinLength || len(s) > TagKeyMaxLength {
		return "", ErrTagKeyInvalidLength
	}

	match, err := regexp.MatchString(TagKeyPattern, s)
	if err != nil {
		return "", err
	}
	if !match {
		return "", ErrTagKeyPatternMismatch
	}

	return TagKey(s), nil
}
