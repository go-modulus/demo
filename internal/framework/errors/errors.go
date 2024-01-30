package errors

import (
	"errors"
	"fmt"
)

func New(message string) error {
	return errors.New(message)
}

func Newf(format string, a ...interface{}) error {
	return New(fmt.Sprintf(format, a...))
}

func NewUserError(code ErrorCode, message string) error {
	err := New(message)
	err = withCode{
		cause: err,
		code:  code,
	}
	err = withMessage{
		cause:   err,
		message: message,
	}
	return err
}

func Wrap(err error, code ErrorCode, message string) error {
	if err == nil {
		return nil
	}
	err = withCode{
		cause: err,
		code:  code,
	}
	err = withMessage{
		cause:   err,
		message: message,
	}
	return err
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}
