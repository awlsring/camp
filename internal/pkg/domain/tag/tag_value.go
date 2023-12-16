package tag

import (
	"fmt"
	"regexp"
)

const (
	TagValueMinLength = 1
	TagValueMaxLength = 128
	TagValuePattern   = "^[a-zA-Z0-9_]+( [a-zA-Z0-9_]+){0,127}$"
)

var (
	ErrTagValueInvalidLength   = fmt.Errorf("tag value must be between 1 and 128 characters")
	ErrTagValuePatternMismatch = fmt.Errorf("tag value must match pattern %s", TagValuePattern)
)

type TagValue string

func (t TagValue) String() string {
	return string(t)
}

func TagValueFromString(s string) (TagValue, error) {
	if len(s) < TagValueMinLength || len(s) > TagValueMaxLength {
		return "", ErrTagValueInvalidLength
	}

	match, err := regexp.MatchString(TagValuePattern, s)
	if err != nil {
		return "", err
	}
	if !match {
		return "", ErrTagValuePatternMismatch
	}

	return TagValue(s), nil
}
