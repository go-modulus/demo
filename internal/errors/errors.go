package errors

import (
	"errors"
	"fmt"
)

type ErrorFlags byte

const (
	ErrorUserFriendly ErrorFlags = 1 << iota
	ErrorDontHandle
)

const (
	DefaultType                = "Error"
	InternalServerErrorCode    = "server.internalError"
	InternalServerErrorMessage = "Something went wrong"
)

type Error struct {
	Type     string
	Code     string
	Message  string
	Metadata map[string]string
	Flags    ErrorFlags
	cause    error
}

func New(code string, message string) *Error {
	return &Error{
		Type:    DefaultType,
		Code:    code,
		Message: message,
	}
}

func Newf(code string, format string, a ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, a...))
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"error: type = %s code = %s message = %s metadata = %v flags = %b cause = %v",
		e.Type,
		e.Code,
		e.Message,
		e.Metadata,
		e.Flags,
		e.cause,
	)
}

func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Type == e.Type && se.Code == e.Code
	}
	return false
}

func (e *Error) Clone() *Error {
	metadata := make(map[string]string, len(e.Metadata))
	for k, v := range e.Metadata {
		metadata[k] = v
	}

	return &Error{
		Type:     e.Type,
		Code:     e.Code,
		Message:  e.Message,
		Metadata: metadata,
		Flags:    e.Flags,
		cause:    e.cause,
	}
}

func (e *Error) WithType(t string) *Error {
	err := e.Clone()
	err.Type = t
	return err
}

func (e *Error) WithCause(cause error) *Error {
	err := e.Clone()
	err.cause = cause
	return err
}

func (e *Error) WithMetadata(md map[string]string) *Error {
	err := e.Clone()
	err.Metadata = md
	return err
}

func (e *Error) WithFlags(flags ErrorFlags) *Error {
	err := e.Clone()
	err.Flags = flags
	return err
}

func (e *Error) ToUserError() *UserError {
	code := e.Code
	message := e.Message

	if e.Flags&ErrorUserFriendly == 0 {
		code = InternalServerErrorCode
		message = InternalServerErrorMessage
	}

	return &UserError{
		Code:       code,
		Message:    message,
		DontHandle: e.Flags&ErrorDontHandle != 0,
		Extra: map[string]interface{}{
			"code": code,
		},
	}
}

func Type(err error) string {
	if err == nil {
		return DefaultType
	}

	return FromError(err).Type
}

func Code(err error) string {
	if err == nil {
		return InternalServerErrorCode
	}

	return FromError(err).Code
}

func Cause(err error) error {
	if err == nil {
		return nil
	}

	return FromError(err).cause
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}

	return New(InternalServerErrorCode, err.Error())
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
