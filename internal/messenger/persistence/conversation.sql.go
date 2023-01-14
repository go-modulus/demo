// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: conversation.sql

package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

const createOrGetConversation = `-- name: CreateOrGetConversation :one
with new_conversation as (
    insert into "messenger"."conversation" (id, sender_id, receiver_id, updated_at, created_at)
    select $1,$2,$3,$4,$5
    where not exists (
        select 1 from "messenger"."conversation"
        where (sender_id = $2 and receiver_id = $3)
        or (receiver_id = $2 and sender_id = $3)
    )
    on conflict do nothing
    returning id, sender_id, receiver_id, updated_at, created_at
) (
    select id, sender_id, receiver_id, updated_at, created_at from "messenger"."conversation"
    where (sender_id = $2 and receiver_id = $3)
    or (receiver_id = $2 and sender_id = $3)

    union all

    select id, sender_id, receiver_id, updated_at, created_at from new_conversation
) limit 1
`

type CreateOrGetConversationParams struct {
	ID         uuid.UUID
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

func (q *Queries) CreateOrGetConversation(ctx context.Context, arg CreateOrGetConversationParams) (MessengerConversation, error) {
	row := q.db.QueryRowContext(ctx, createOrGetConversation,
		arg.ID,
		arg.SenderID,
		arg.ReceiverID,
		arg.UpdatedAt,
		arg.CreatedAt,
	)
	var i MessengerConversation
	err := row.Scan(
		&i.ID,
		&i.SenderID,
		&i.ReceiverID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getConversation = `-- name: GetConversation :one
select id, sender_id, receiver_id, updated_at, created_at from "messenger"."conversation" where id = $1 limit 1
`

func (q *Queries) GetConversation(ctx context.Context, id uuid.UUID) (MessengerConversation, error) {
	row := q.db.QueryRowContext(ctx, getConversation, id)
	var i MessengerConversation
	err := row.Scan(
		&i.ID,
		&i.SenderID,
		&i.ReceiverID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const paginateMyConversations = `-- name: PaginateMyConversations :many
select id, sender_id, receiver_id, updated_at, created_at from "messenger"."conversation"
where
    (sender_id = $1 or receiver_id = $1)
    and (
        $2::timestamp is null
        or updated_at < $2
        or ($3::uuid is null or (updated_at = $2 and id > $3))
    )
order by updated_at desc, id
limit $4
`

type PaginateMyConversationsParams struct {
	ViewerID       uuid.UUID
	AfterUpdatedAt sql.NullTime
	AfterID        uuid.NullUUID
	First          int32
}

func (q *Queries) PaginateMyConversations(ctx context.Context, arg PaginateMyConversationsParams) ([]MessengerConversation, error) {
	rows, err := q.db.QueryContext(ctx, paginateMyConversations,
		arg.ViewerID,
		arg.AfterUpdatedAt,
		arg.AfterID,
		arg.First,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MessengerConversation
	for rows.Next() {
		var i MessengerConversation
		if err := rows.Scan(
			&i.ID,
			&i.SenderID,
			&i.ReceiverID,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
