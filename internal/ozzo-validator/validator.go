package validator

import (
	"boilerplate/internal/framework"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func AsAppValidationErrors(
	err error,
) *framework.ValidationErrors {
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

		return framework.NewValidationErrors(errors)
	}
	return framework.NewValidationErrors(
		[]framework.ValidationError{
			{
				Field:      "",
				Identifier: framework.UnknownError,
				Err:        err.Error(),
			},
		},
	)
}
