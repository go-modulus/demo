package errors

import "errors"

var (
	ErrNonRetryable = errors.New("non-retryable error")
)
