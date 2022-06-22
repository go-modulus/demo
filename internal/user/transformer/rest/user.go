package transformer

import (
	"demo/internal/user/action"
	"demo/internal/user/storage"
)

func TransformUser(user storage.User) *action.User {
	return &action.User{
		Id:   user.ID.String(),
		Name: user.Name,
	}
}
