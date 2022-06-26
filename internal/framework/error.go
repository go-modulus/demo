package framework

import "context"

type ErrorIdentifier string

const WrongRequestDecoding ErrorIdentifier = "WrongRequestDecoding"
const InvalidRequest ErrorIdentifier = "InvalidRequest"
const UnprocessableEntity ErrorIdentifier = "UnprocessableEntity"
const UnknownError ErrorIdentifier = "UnknownError"

type CommonError struct {
	Identifier ErrorIdentifier
	Err        string
}

func (e *CommonError) Error() string {
	return e.Err
}

func NewCommonError(identifier ErrorIdentifier, err string) *CommonError {
	return &CommonError{
		Identifier: identifier,
		Err:        err,
	}
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
		identifier = commonErr.Identifier
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
