package errors

import (
	"errors"

	"golang.org/x/text/message"
)

func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return withMessage{
		cause:   err,
		message: message,
	}
}

type withMessage struct {
	cause   error
	message string
}

func (w withMessage) Error() string   { return w.cause.Error() }
func (w withMessage) Cause() error    { return w.cause }
func (w withMessage) Message() string { return w.message }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w withMessage) Unwrap() error { return w.cause }

func Message(t *message.Printer, err error) string {
	type withMessage interface {
		Message() string
	}
	var wm withMessage
	if errors.As(err, &wm) {
		return wm.Message()
	}

	type withMessagePrinter interface {
		Message(*message.Printer) string
	}
	var wmp withMessagePrinter
	if errors.As(err, &wmp) {
		return wmp.Message(t)
	}

	return "Internal Server Error"
}
