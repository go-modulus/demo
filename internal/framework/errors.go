package framework

import (
	"context"
	"fmt"
	"go.uber.org/fx"
)

type Error struct {
	Message string
	Code    string
}

func (e *Error) Error() string {
	return e.Message
}

type BusinessLogicError struct {
	Error
}

type ErrorFilter func(ctx context.Context, err error) bool
type ErrorListener func(ctx context.Context, err error) error

type ErrorHandler struct {
	filters   []ErrorFilter
	listeners []ErrorListener
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
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

func ErrorsModule() fx.Option {
	return fx.Module(
		"errors",
		fx.Provide(NewErrorHandler),
	)
}
