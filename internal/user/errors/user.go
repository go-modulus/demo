package errors

import "demo/internal/errors"

const UserNotFoundCode errors.ErrorCode = "user.notFound"

func NewUserNotFound(id string) *errors.Error {
	err := errors.NewNotFoundError(
		UserNotFoundCode,
		"This user not found",
	)

	return err.WithMetadata(map[string]string{"id": id})
}
