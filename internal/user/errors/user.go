package errors

import (
	"demo/internal/errors"
	"fmt"
)

const UserNotFoundCode errors.ErrorCode = "user.notFound"

func NewUserNotFound(id string) *errors.Error {
	err := errors.NewNotFoundError(
		UserNotFoundCode,
		fmt.Sprintf("user with id %s not found", id),
	)

	return err.WithMetadata(map[string]string{"id": id})
}
