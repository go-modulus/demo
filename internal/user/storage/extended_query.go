package storage

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/gofrs/uuid"
)

type UuidSlice []uuid.UUID

func (j UuidSlice) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *UuidSlice) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return nil
	}
	return nil
}

func (q *Queries) GetUsersMap(ctx context.Context, ids []uuid.UUID) (map[string]User, error) {
	users, err := q.GetUsersByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	res := make(map[string]User, len(users))
	for _, user := range users {
		res[user.ID.String()] = user
	}

	return res, nil
}
