package errors

const (
	BusinessLogicError = "BusinessLogicError"
	NotFoundError      = "NotFoundError"
)

func NewBusinessLogicError(code string, message string) *Error {
	return New(code, message).
		WithType(BusinessLogicError).
		WithFlags(ErrorUserFriendly)
}

func IsBusinessLogicError(err error) bool {
	return Type(err) == BusinessLogicError
}

func NewNotFoundError(code string, message string) *Error {
	return New(code, message).
		WithType(NotFoundError).
		WithFlags(ErrorUserFriendly | ErrorDontHandle)
}

func IsNotFoundError(err error) bool {
	return Type(err) == NotFoundError
}
