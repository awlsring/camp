package handler

import (
	camperror "github.com/awlsring/camp/internal/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"errors"
)

func grpcError(err error) error {
	var campErr *camperror.Error
	switch {
	case errors.As(err, &campErr):
		e := campErr.CampError()
		switch e {
		case camperror.ErrInternalFailure:
			return status.Errorf(codes.Internal, err.Error())
		case camperror.ErrResourceNotFound:
			return status.Errorf(codes.NotFound, err.Error())
		case camperror.ErrUnathorized:
			return status.Errorf(codes.Unauthenticated, err.Error())
		case camperror.ErrValidation:
			return status.Errorf(codes.InvalidArgument, err.Error())
		case camperror.ErrDuplicate:
			return status.Errorf(codes.AlreadyExists, err.Error())
		}
	}
	return status.Errorf(codes.Internal, err.Error())
}
