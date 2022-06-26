package errors

import (
	"context"
	"errors"
	"fmt"
)

type ErrorFilter func(ctx context.Context, err error) bool
type ErrorListener func(ctx context.Context, err error) error

type ErrorHandler struct {
	filters   []ErrorFilter
	listeners []ErrorListener
}

func NewErrorHandler() *ErrorHandler {
	eh := &ErrorHandler{}

	eh.AttachFilter(func(ctx context.Context, err error) bool {
		is := errors.Is(
			context.Canceled,
			err,
		)
		if is {
			return false
		}

		return false == errors.Is(
			context.DeadlineExceeded,
			err,
		)
	})

	return eh
}

func (h *ErrorHandler) AttachFilter(filter ErrorFilter) {
	h.filters = append(h.filters, filter)
}

func (h *ErrorHandler) AttachListener(listener ErrorListener) {
	h.listeners = append(h.listeners, listener)
}

func (h *ErrorHandler) Handle(ctx context.Context, err error) {
	for _, filter := range h.filters {
		shouldContinue := filter(ctx, err)

		if shouldContinue == false {
			return
		}
	}

	for _, listener := range h.listeners {
		listenerErr := listener(ctx, err)

		if listenerErr != nil {
			panic(fmt.Errorf("cannot handle error: %w", err))
		}
	}
}
