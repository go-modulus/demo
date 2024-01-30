package errors

import (
	"errors"
)

type WithStack interface {
	Stack() string
}

func Loggable(err error) bool {
	type withLoggable interface {
		IsLoggable() bool
	}
	var wl withLoggable
	if errors.As(err, &wl) {
		return wl.IsLoggable()
	}

	return false
}

func Stack(err error) string {
	var ws WithStack
	if errors.As(err, &ws) {
		return ws.Stack()
	}

	return ""
}
