package errors

import (
	"context"
	"fmt"
	application "github.com/debugger84/modulus-application"
)

const cannotUpdateUser application.ErrorIdentifier = "CannotUpdateUser"

func CannotUpdateUser(ctx context.Context, id string) *application.ActionResponse {
	return &application.ActionResponse{
		StatusCode: 422,
		Error: &application.ActionError{
			Ctx:              ctx,
			Identifier:       cannotUpdateUser,
			Err:              fmt.Errorf("User with id %s cannot be updated", id),
			ValidationErrors: nil,
		},
	}
}
