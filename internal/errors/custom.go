package errors

const (
	BusinessLogicError ErrorType = "BusinessLogicError"
	BadRequestError    ErrorType = "BadRequestError"
	NotFoundError      ErrorType = "NotFoundError"
	UnauthorizedError  ErrorType = "UnauthorizedError"
)

func NewBusinessLogicError(code ErrorCode, message string) *Error {
	return New(code, message).
		WithType(BusinessLogicError).
		WithFlags(ErrorUserFriendly)
}

func IsBusinessLogicError(err error) bool {
	return Type(err) == BusinessLogicError
}

func NewBadRequestError(code ErrorCode, message string) *Error {
	return New(code, message).
		WithType(BadRequestError).
		WithFlags(ErrorUserFriendly | ErrorDontHandle)
}

func IsBadRequestError(err error) bool {
	return Type(err) == BadRequestError
}

func NewNotFoundError(code ErrorCode, message string) *Error {
	return New(code, message).
		WithType(NotFoundError).
		WithFlags(ErrorUserFriendly | ErrorDontHandle)
}

func IsNotFoundError(err error) bool {
	return Type(err) == NotFoundError
}

func NewUnauthorizedError(code ErrorCode, message string) *Error {
	return New(code, message).
		WithType(UnauthorizedError).
		WithFlags(ErrorUserFriendly | ErrorDontHandle)
}

func IsUnauthorizedError(err error) bool {
	return Type(err) == UnauthorizedError
}
