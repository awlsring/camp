package machine

import (
	"errors"

	camperror "github.com/awlsring/camp/internal/pkg/errors"
)

var (
	ErrEmptyValue = errors.New("empty value given, must not be null")
)

func NonNullString[T ~string](v string) (T, error) {
	if v == "" {
		return "", camperror.New(camperror.ErrValidation, ErrEmptyValue)
	}
	return T(v), nil
}
