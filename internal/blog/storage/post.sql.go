// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: post.sql

package storage

import (
	"context"

	uuid "github.com/gofrs/uuid"
)

const countPosts = `-- name: CountPosts :one
select count(*) as count
from blog."post" as p
`

func (q *Queries) CountPosts(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countPosts)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPost = `-- name: CreatePost :one
insert into blog."post" (id, title, body, author_id, slug)
values ($1, $2, $3, $4, $5)
RETURNING id, title, body, author_id, slug, status, created_at, published_at, updated_at
`

type CreatePostParams struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Title    string    `db:"title" json:"title"`
	Body     string    `db:"body" json:"body"`
	AuthorID uuid.UUID `db:"author_id" json:"authorId"`
	Slug     string    `db:"slug" json:"slug"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRow(ctx, createPost,
		arg.ID,
		arg.Title,
		arg.Body,
		arg.AuthorID,
		arg.Slug,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Body,
		&i.AuthorID,
		&i.Slug,
		&i.Status,
		&i.CreatedAt,
		&i.PublishedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
delete
from blog."post"
where id = $1::uuid
`

func (q *Queries) DeletePost(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deletePost, id)
	return err
}

const getPost = `-- name: GetPost :one
select id, title, body, author_id, slug, status, created_at, published_at, updated_at
from blog."post"
where id = $1::uuid
LIMIT 1
`

func (q *Queries) GetPost(ctx context.Context, id uuid.UUID) (Post, error) {
	row := q.db.QueryRow(ctx, getPost, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Body,
		&i.AuthorID,
		&i.Slug,
		&i.Status,
		&i.CreatedAt,
		&i.PublishedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listPosts = `-- name: ListPosts :many
select p.id, p.title, p.body, p.author_id, p.slug, p.status, p.created_at, p.published_at, p.updated_at
from blog."post" as p
order by p.published_at desc
limit $2 offset $1
`

type ListPostsParams struct {
	After int32 `db:"after" json:"after"`
	Count int32 `db:"count" json:"count"`
}

func (q *Queries) ListPosts(ctx context.Context, arg ListPostsParams) ([]Post, error) {
	rows, err := q.db.Query(ctx, listPosts, arg.After, arg.Count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Body,
			&i.AuthorID,
			&i.Slug,
			&i.Status,
			&i.CreatedAt,
			&i.PublishedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const publishPost = `-- name: PublishPost :one
update blog."post"
set published_at = now(),
status = 'published'
where id = $1::uuid
RETURNING id, title, body, author_id, slug, status, created_at, published_at, updated_at
`

func (q *Queries) PublishPost(ctx context.Context, id uuid.UUID) (Post, error) {
	row := q.db.QueryRow(ctx, publishPost, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Body,
		&i.AuthorID,
		&i.Slug,
		&i.Status,
		&i.CreatedAt,
		&i.PublishedAt,
		&i.UpdatedAt,
	)
	return i, err
}
