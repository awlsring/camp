package errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcError(err error) error {
	var campErr *Error
	switch {
	case errors.As(err, &campErr):
		e := campErr.CampError()
		switch e {
		case ErrInternalFailure:
			return status.Errorf(codes.Internal, err.Error())
		case ErrResourceNotFound:
			return status.Errorf(codes.NotFound, err.Error())
		case ErrUnathorized:
			return status.Errorf(codes.Unauthenticated, err.Error())
		case ErrValidation:
			return status.Errorf(codes.InvalidArgument, err.Error())
		case ErrDuplicate:
			return status.Errorf(codes.AlreadyExists, err.Error())
		}
	}
	return status.Errorf(codes.Internal, err.Error())
}
