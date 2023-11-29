package exception

type ExceptionType int

const (
	ExceptionTypeInternalServerException ExceptionType = iota
	ExceptionTypeUnauthorizedException
	ExceptionTypeInvalidInputException
	ExceptionTypeResourceNotFoundException
	ExceptionTypeInternalFailureException
	ExceptionTypeSerializationException
	ExceptionTypeUnknownOperationException
)

func (e ExceptionType) String() string {
	switch e {
	case ExceptionTypeInternalServerException:
		return "InternalServerException"
	case ExceptionTypeUnauthorizedException:
		return "UnauthorizedException"
	case ExceptionTypeInvalidInputException:
		return "InvalidInputException"
	case ExceptionTypeResourceNotFoundException:
		return "ResourceNotFoundException"
	case ExceptionTypeSerializationException:
		return "SerializationException"
	case ExceptionTypeUnknownOperationException:
		return "UnknownOperationException"
	default:
		return "InternalFailureException"
	}
}
