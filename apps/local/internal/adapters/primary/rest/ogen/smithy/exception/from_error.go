package exception

import (
	"errors"

	camperror "github.com/awlsring/camp/internal/pkg/errors"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
)

func FromError(err error) *Exception {
	var (
		ctError *validate.InvalidContentTypeError
		ogenErr ogenerrors.Error
		campErr *camperror.Error
	)
	switch {
	case errors.Is(err, ht.ErrNotImplemented):
		return NewUnknownOperationException(err.Error())
	case errors.As(err, &ctError):
		return NewInternalServerException(err.Error())
	case errors.As(err, &ogenErr):
		code := ogenErr.Code()
		switch code {
		case 400:
			return NewValidationError(err.Error())
		}
	case errors.As(err, &campErr):
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

	return NewInternalServerException(err.Error())
}
