package framework

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type ErrorIdentifier string

const WrongRequestDecoding ErrorIdentifier = "WrongRequestDecoding"
const InvalidRequest ErrorIdentifier = "InvalidRequest"
const UnprocessableEntity ErrorIdentifier = "UnprocessableEntity"
const UnknownError ErrorIdentifier = "UnknownError"

var p = message.NewPrinter(language.English)

type UserError interface {
	Identifier() ErrorIdentifier
	Message(*message.Printer) string
}

type ExtraError interface {
	Extra() map[string]any
}

type CommonError struct {
	identifier   ErrorIdentifier
	errTpl       string
	tplVariables []any
}

func (e *CommonError) Identifier() ErrorIdentifier {
	return e.identifier
}

func (e *CommonError) Message(printer *message.Printer) string {
	return printer.Sprintf(e.errTpl, e.tplVariables...)
}

func (e *CommonError) Error() string {
	return fmt.Sprintf(e.errTpl, e.tplVariables...)
}

func NewTranslatedError(ctx context.Context, identifier ErrorIdentifier, errTpl string, variables ...any) *CommonError {
	t := GetTranslator(ctx)
	msg := t.Sprintf(errTpl, variables...)
	return &CommonError{
		identifier:   identifier,
		errTpl:       msg,
		tplVariables: variables,
	}
}

func NewCommonError(identifier ErrorIdentifier, errTpl string, variables ...any) *CommonError {
	// it is a hack to mark the error for extracting to the translation file
	_ = p.Sprintf(errTpl, variables...)
	return &CommonError{
		identifier:   identifier,
		errTpl:       errTpl,
		tplVariables: variables,
	}
}

func (e *CommonError) WithTplVariables(variables ...any) error {
	return NewCommonError(e.identifier, e.errTpl, variables...)
}

func (e *CommonError) Is(err error) bool {
	if cErr, ok := err.(*CommonError); ok {
		return cErr.identifier == e.identifier
	}
	return false
}

type ActionError struct {
	Ctx              context.Context
	Identifier       ErrorIdentifier
	Err              error
	ValidationErrors []ValidationError
}

func (e *ActionError) Error() string {
	return e.Err.Error()
}

type ValidationError struct {
	Field      string
	Identifier ErrorIdentifier
	Err        string
}

func NewValidationError(field string, err string, identifier ErrorIdentifier) *ValidationError {
	return &ValidationError{Field: field, Err: err, Identifier: identifier}
}

func (e ValidationError) Error() string {
	return e.Err
}

func (e ValidationError) Is(err error) bool {
	if cErr, ok := err.(*CommonError); ok {
		return cErr.identifier == e.Identifier
	}
	if cErr, ok := err.(*ValidationError); ok {
		return cErr.Identifier == e.Identifier
	}
	return false
}

func (e ValidationError) Message(printer *message.Printer) string {
	return printer.Sprintf(e.Err)
}

type ValidationErrors struct {
	errors []ValidationError
}

func NewValidationErrors(errors []ValidationError) *ValidationErrors {
	if len(errors) == 0 {
		errors = []ValidationError{
			{
				Field:      "",
				Identifier: UnknownError,
				Err:        "Unknown error",
			},
		}
	}
	return &ValidationErrors{errors: errors}
}

func (e ValidationErrors) Error() string {
	return e.errors[0].Error()
}

func (e ValidationErrors) Errors() []ValidationError {
	return e.errors
}

func (e ValidationErrors) ErrorMessages() map[string]string {
	messages := make(map[string]string)
	for _, validationError := range e.errors {
		messages[validationError.Field] = validationError.Message(p)
	}
	return messages
}

func (e ValidationErrors) Is(err error) bool {
	if cErr, ok := err.(*CommonError); ok {
		for _, validationError := range e.errors {
			if cErr.identifier == validationError.Identifier {
				return true
			}
		}
		return false
	}
	if cErr, ok := err.(*ValidationError); ok {
		for _, validationError := range e.errors {
			if cErr.Identifier == validationError.Identifier {
				return true
			}
		}
		return false
	}
	return false
}

func NewServerErrorResponse(ctx context.Context, identifier ErrorIdentifier, err error) ActionResponse {
	return ActionResponse{
		StatusCode: 500,
		Error: &ActionError{
			Ctx:              ctx,
			Identifier:       identifier,
			Err:              err,
			ValidationErrors: nil,
		},
	}
}

func NewUnprocessableEntityResponse(ctx context.Context, err error) ActionResponse {
	code := 500
	identifier := UnprocessableEntity
	if commonErr, ok := err.(*CommonError); ok {
		code = 422
		identifier = commonErr.identifier
	}
	return ActionResponse{
		StatusCode: code,
		Error: &ActionError{
			Ctx:              ctx,
			Identifier:       identifier,
			Err:              err,
			ValidationErrors: nil,
		},
	}
}

func NewValidationErrorResponse(ctx context.Context, errors []ValidationError) ActionResponse {
	if len(errors) == 0 {
		errors[0] = ValidationError{
			Err: "Unknown error",
		}
	}
	return ActionResponse{
		StatusCode: 400,
		Error: &ActionError{
			Ctx:              ctx,
			Identifier:       InvalidRequest,
			Err:              errors[0],
			ValidationErrors: errors,
		},
	}
}
