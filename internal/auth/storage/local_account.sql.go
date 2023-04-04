// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: local_account.sql

package storage

import (
	"context"

	uuid "github.com/gofrs/uuid"
	null "gopkg.in/guregu/null.v4"
)

const createLocalAccount = `-- name: CreateLocalAccount :one
INSERT INTO auth.local_account (user_id, email, nickname, phone, password, created_at)
VALUES ($1::uuid, $2, $3,
        $4,
        $5::text, NOW())
 RETURNING user_id, email, nickname, phone, password, created_at
`

type CreateLocalAccountParams struct {
	UserID       uuid.UUID   `db:"user_id" json:"userID"`
	Email        null.String `db:"email" json:"email"`
	Nickname     null.String `db:"nickname" json:"nickname"`
	Phone        null.String `db:"phone" json:"phone"`
	PasswordHash string      `db:"password_hash" json:"passwordHash"`
}

func (q *Queries) CreateLocalAccount(ctx context.Context, arg CreateLocalAccountParams) (LocalAccount, error) {
	row := q.db.QueryRow(ctx, createLocalAccount,
		arg.UserID,
		arg.Email,
		arg.Nickname,
		arg.Phone,
		arg.PasswordHash,
	)
	var i LocalAccount
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Nickname,
		&i.Phone,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const deleteLocalAccount = `-- name: DeleteLocalAccount :execrows
delete from auth.local_account where user_id = $1::uuid
`

func (q *Queries) DeleteLocalAccount(ctx context.Context, userID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteLocalAccount, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
