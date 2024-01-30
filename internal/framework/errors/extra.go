package errors

import (
	"errors"
	"fmt"
)

func WithExtra(err error, extra map[string]any) error {
	if err == nil {
		return nil
	}
	return withExtra{
		cause: err,
		extra: extra,
	}
}

type withExtra struct {
	cause error
	extra map[string]any
}

func (w withExtra) Error() string {
	return fmt.Sprintf("%v: %v", w.cause, w.extra)
}
func (w withExtra) Cause() error          { return w.cause }
func (w withExtra) Extra() map[string]any { return w.extra }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w withExtra) Unwrap() error { return w.cause }

func Extra(err error) map[string]any {
	type withExtra interface {
		Extra() map[string]any
	}
	var we withExtra
	if errors.As(err, &we) {
		return we.Extra()
	}

	return nil
}
