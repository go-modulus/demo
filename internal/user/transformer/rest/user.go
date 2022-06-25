package transformer

import (
	"boilerplate/internal/user/action"
	"boilerplate/internal/user/storage"
)

func TransformUser(user storage.User) *action.User {
	return &action.User{
		Id:   user.ID.String(),
		Name: user.Name,
	}
}
