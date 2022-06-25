package errors

import (
	"context"
	"errors"
	"fmt"
	application "github.com/debugger84/modulus-application"
)

const userNotFound application.ErrorIdentifier = "UserNotFound"
const cannotUpdateUser application.ErrorIdentifier = "CannotUpdateUser"

func UserNotFound(ctx context.Context, id string) application.ActionResponse {
	return application.ActionResponse{
		StatusCode: 404,
		Error: &application.ActionError{
			Ctx:              ctx,
			Identifier:       userNotFound,
			Err:              errors.New(fmt.Sprintf("User with id %s is not found", id)),
			ValidationErrors: nil,
		},
	}
}

func CannotUpdateUser(ctx context.Context, id string) application.ActionResponse {
	return application.ActionResponse{
		StatusCode: 422,
		Error: &application.ActionError{
			Ctx:              ctx,
			Identifier:       cannotUpdateUser,
			Err:              errors.New(fmt.Sprintf("User with id %s cannot be updated", id)),
			ValidationErrors: nil,
		},
	}
}
