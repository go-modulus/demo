package validator

import "context"

type Validatable interface {
	Validate(ctx context.Context) *ValidationError
}
