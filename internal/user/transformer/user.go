package transformer

import (
	"boilerplate/internal/graph/model"
	"boilerplate/internal/user/storage"
	"encoding/base64"
	"encoding/json"
	"time"
)

type UsersListCursor struct {
	RegisteredAt time.Time `json:"ra"`
	Id           string    `json:"id"`
}

func NewUsersListCursor(user storage.User) *UsersListCursor {
	return &UsersListCursor{
		RegisteredAt: user.RegisteredAt,
		Id:           user.ID.String(),
	}
}

func NewUsersListCursorFromString(cursor string) *UsersListCursor {
	var res *UsersListCursor
	buf, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return nil
	}
	return res
}

func (c *UsersListCursor) ToString() string {
	val, err := json.Marshal(&c)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(val)
}

func TransformUser(user storage.User) *model.User {
	return &model.User{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}
}

func TransformUserList(
	users []storage.User,
	maxCount int,
	cursorFactory func(user storage.User) string,
) *model.UserList {
	result := model.UserList{
		Edges:       make([]*model.UserEdge, 0, maxCount),
		HasNextPage: maxCount < len(users),
	}
	for i, user := range users {
		if i == maxCount {
			break
		}
		edge := model.UserEdge{
			Cursor: cursorFactory(user),
			Node:   TransformUser(user),
		}
		result.Edges = append(result.Edges, &edge)
		result.EndCursor = edge.Cursor
	}
	return &result
}
