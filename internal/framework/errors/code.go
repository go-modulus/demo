package errors

import (
	"ecs/internal/infra/framework"
	"errors"
)

type ErrorCode string

func (c ErrorCode) String() string { return string(c) }

const (
	InternalServerErrorCode ErrorCode = "server.internalError"
)

func WithCode(err error, code ErrorCode) error {
	return withCode{
		cause: err,
		code:  code,
	}
}

type withCode struct {
	cause error
	code  ErrorCode
}

func (w withCode) Error() string   { return "<" + w.code.String() + "> " + w.cause.Error() }
func (w withCode) Cause() error    { return w.cause }
func (w withCode) Code() ErrorCode { return w.code }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w withCode) Unwrap() error { return w.cause }

func Code(err error) ErrorCode {
	type withCode interface {
		Code() ErrorCode
	}
	var wc withCode
	if errors.As(err, &wc) {
		return wc.Code()
	}

	// fallback to legacy identifier
	type withIdentifier interface {
		Identifier() framework.ErrorIdentifier
	}
	var wi withIdentifier
	if errors.As(err, &wi) {
		return ErrorCode(wi.Identifier())
	}

	return InternalServerErrorCode
}
