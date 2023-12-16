package values

import "errors"

var (
	ErrEmptyValue = errors.New("empty value given, must not be null")
)

func NonNullString[T ~string](v string) (T, error) {
	if v == "" {
		return "", ErrEmptyValue
	}
	return T(v), nil
}
