package exception

type Exception struct {
	HttpCode      int
	ExceptionType ExceptionType
	Message       string
}

func (e *Exception) Error() string {
	return e.Message
}

func (e *Exception) Code() int {
	return e.HttpCode
}

func (e *Exception) Type() ExceptionType {
	return e.ExceptionType
}

func (e *Exception) AsJsonMessage() []byte {
	return []byte(`{"message": "` + e.Message + `"}`)
}

func NewInternalServerException(msg string) *Exception {
	return &Exception{
		HttpCode:      500,
		ExceptionType: ExceptionTypeInternalServerException,
		Message:       msg,
	}
}

func NewUnauthorizedError(msg string) *Exception {
	return &Exception{
		HttpCode:      401,
		ExceptionType: ExceptionTypeUnauthorizedException,
		Message:       msg,
	}
}

func NewInvalidInputError(msg string) *Exception {
	return &Exception{
		HttpCode:      400,
		ExceptionType: ExceptionTypeInvalidInputException,
		Message:       msg,
	}
}

func NewResourceNotFoundError(msg string) *Exception {
	return &Exception{
		HttpCode:      404,
		ExceptionType: ExceptionTypeResourceNotFoundException,
		Message:       msg,
	}
}

func NewInternalFailureException(msg string) *Exception {
	return &Exception{
		HttpCode:      500,
		ExceptionType: ExceptionTypeInternalFailureException,
		Message:       msg,
	}
}

func NewSerializationException(msg string) *Exception {
	return &Exception{
		HttpCode:      400,
		ExceptionType: ExceptionTypeSerializationException,
		Message:       msg,
	}
}

func NewUnknownOperationException(msg string) *Exception {
	return &Exception{
		HttpCode:      400,
		ExceptionType: ExceptionTypeUnknownOperationException,
		Message:       msg,
	}
}
