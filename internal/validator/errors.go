package validator

import (
	"demo/internal/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

type FieldValidationError struct {
	Field   string
	Code    string
	Message string
}

func (f FieldValidationError) ToUserError() map[string]interface{} {
	return map[string]interface{}{
		"field":   f.Field,
		"code":    f.Code,
		"message": f.Message,
	}
}

type ValidationError struct {
	BaseError *errors.Error
	Fields    []FieldValidationError
}

func NewValidationError(fields []FieldValidationError) *ValidationError {
	if len(fields) == 0 {
		return nil
	}

	return &ValidationError{
		BaseError: errors.
			New("validation.failed", "Validation failed").
			WithType("ValidationError").
			WithFlags(errors.ErrorUserFriendly | errors.ErrorDontHandle),
		Fields: fields,
	}
}

func FromOzzoError(err error) *ValidationError {
	if err == nil {
		return nil
	}

	ozzoErrors, ok := err.(validation.Errors)
	if !ok {
		return nil
	}

	fieldErrors := make([]FieldValidationError, 0, len(ozzoErrors))
	for key, val := range ozzoErrors {
		errObj, ok := val.(validation.ErrorObject)
		if !ok {
			continue
		}

		fieldErrors = append(
			fieldErrors,
			FieldValidationError{
				Field:   key,
				Code:    strings.Replace(errObj.Code(), "_", ".", 1),
				Message: errObj.Error(),
			},
		)
	}

	return NewValidationError(fieldErrors)
}
func (v ValidationError) Error() string {
	return v.BaseError.Error()
}

func (v ValidationError) ToUserError() *errors.UserError {
	userError := v.BaseError.ToUserError()

	fields := make([]map[string]interface{}, 0, len(v.Fields))
	for _, field := range v.Fields {
		fields = append(
			fields,
			field.ToUserError(),
		)
	}

	userError.Extra["fields"] = fields

	return userError
}
