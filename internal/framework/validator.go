package framework

import (
	"context"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
)

type StructValidator interface {
	ValidateStruct(obj any) []ValidationError
}

type VarValidator interface {
	ValidateVar(variable any, rule string) *ValidationError
}

type ValidatableStruct interface {
	Validate(ctx context.Context) []ValidationError
}

type DefaultValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func NewDefaultValidator(logger Logger) StructValidator {
	uni := ut.New(en.New())
	translator, _ := uni.GetTranslator("en")
	validate := validator.New()
	err := enTranslations.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		logger.Error(context.Background(), "Cannot register default translations for validator")
	}
	validate.RegisterTagNameFunc(
		func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		},
	)
	return &DefaultValidator{validator: validate, translator: translator}
}

func (v *DefaultValidator) ValidateStruct(obj any) []ValidationError {
	err := v.validator.Struct(obj)
	if err != nil {
		if validatorErr, ok := err.(validator.ValidationErrors); ok {
			result := make([]ValidationError, len(validatorErr))
			for i, validationError := range validatorErr {
				result[i] = *NewValidationError(
					validationError.Field(),
					validationError.Translate(v.translator),
					ErrorIdentifier(validationError.Error()),
				)
			}
			return result
		} else {
			return []ValidationError{
				*NewValidationError(
					"",
					err.Error(),
					InvalidRequest,
				),
			}
		}
	}
	return nil
}

func (v *DefaultValidator) ValidateVar(variable any, rule string) *ValidationError {
	err := v.validator.Var(variable, rule)
	if err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			return NewValidationError(
				validationErr[0].Field(),
				validationErr[0].Translate(v.translator),
				ErrorIdentifier(validationErr[0].Error()),
			)
		} else {
			return NewValidationError(
				"",
				err.Error(),
				InvalidRequest,
			)
		}
	}
	return nil
}
