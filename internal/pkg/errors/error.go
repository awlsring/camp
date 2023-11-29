package errors

import "fmt"

var (
	ErrInternalFailure  = fmt.Errorf("internal failure")
	ErrResourceNotFound = fmt.Errorf("resource not found")
	ErrUnathorized      = fmt.Errorf("unauthorized")
	ErrValidation       = fmt.Errorf("validation error")
)

type Error struct {
	camp error
	err  error
}

func New(camp, err error) *Error {
	return &Error{
		camp: camp,
		err:  err,
	}
}

func (e *Error) Error() string {
	return fmt.Errorf("%w: %w", e.camp, e.err).Error()
}

func (e *Error) CampError() error {
	return e.camp
}

func (e *Error) WrapperError() error {
	return e.err
}
