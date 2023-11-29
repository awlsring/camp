package exception

import (
	"errors"

	camperror "github.com/awlsring/camp/internal/pkg/errors"
)

func FromError(e error) *Exception {
	var campErr *camperror.Error
	if errors.As(e, &campErr) {
		e := campErr.CampError()
		switch e {
		case camperror.ErrInternalFailure:
			return NewInternalServerException(e.Error())
		case camperror.ErrResourceNotFound:
			return NewResourceNotFoundError(e.Error())
		case camperror.ErrUnathorized:
			return NewUnauthorizedError(e.Error())
		case camperror.ErrValidation:
			return NewInvalidInputError(e.Error())
		}
	}
	return NewInternalServerException(e.Error())
}
