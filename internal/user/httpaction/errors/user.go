package errors

import (
	"boilerplate/internal/framework"
	"context"
	"errors"
	"fmt"
)

const userNotFound framework.ErrorIdentifier = "UserNotFound"
const cannotUpdateUser framework.ErrorIdentifier = "CannotUpdateUser"

func UserNotFound(ctx context.Context, id string) framework.ActionResponse {
	return framework.ActionResponse{
		StatusCode: 404,
		Error: &framework.ActionError{
			Ctx:              ctx,
			Identifier:       userNotFound,
			Err:              errors.New(fmt.Sprintf("User with id %s is not found", id)),
			ValidationErrors: nil,
		},
	}
}

func CannotUpdateUser(ctx context.Context, id string) framework.ActionResponse {
	return framework.ActionResponse{
		StatusCode: 422,
		Error: &framework.ActionError{
			Ctx:              ctx,
			Identifier:       cannotUpdateUser,
			Err:              errors.New(fmt.Sprintf("User with id %s cannot be updated", id)),
			ValidationErrors: nil,
		},
	}
}
