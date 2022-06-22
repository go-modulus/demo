package validator

import (
	"demo/internal/framework"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func AsAppValidationErrors(
	err error,
) []framework.ValidationError {
	if err == nil {
		return nil
	}

	if errorSet, ok := err.(validation.Errors); ok {
		errors := make([]framework.ValidationError, 0, len(errorSet))
		for key, val := range errorSet {
			if errorSet, ok := val.(validation.ErrorObject); ok {
				errors = append(
					errors, framework.ValidationError{
						Field:      key,
						Err:        errorSet.Error(),
						Identifier: framework.ErrorIdentifier(errorSet.Code()),
					},
				)
			}
		}

		return errors
	}
	return []framework.ValidationError{
		{
			Field:      "",
			Identifier: framework.UnknownError,
			Err:        err.Error(),
		},
	}
}
