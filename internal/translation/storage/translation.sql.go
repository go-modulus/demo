// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: translation.sql

package storage

import (
	"context"

	"github.com/jackc/pgtype"
)

const deleteTransactionsByKey = `-- name: DeleteTransactionsByKey :exec
delete from "translation"."translation"
where key = $1
`

func (q *Queries) DeleteTransactionsByKey(ctx context.Context, key string) error {
	_, err := q.db.Exec(ctx, deleteTransactionsByKey, key)
	return err
}

const findTranslationValue = `-- name: FindTranslationValue :one
select value
from "translation"."translation"
where key = $1
  and path = $2
  and locale = $3
`

type FindTranslationValueParams struct {
	Key    string `db:"key" json:"key"`
	Path   Path   `db:"path" json:"path"`
	Locale Locale `db:"locale" json:"locale"`
}

func (q *Queries) FindTranslationValue(ctx context.Context, arg FindTranslationValueParams) (string, error) {
	row := q.db.QueryRow(ctx, findTranslationValue, arg.Key, arg.Path, arg.Locale)
	var value string
	err := row.Scan(&value)
	return value, err
}

const findTranslations = `-- name: FindTranslations :many
SELECT t.key, t.path, t.locale, t.value, t.created_at, t.updated_at, t.pseudo_id
FROM jsonb_to_recordset($1::jsonb) AS input_ids(key text, path translation.path, locale translation.locale)
JOIN "translation"."translation" as t
     on input_ids.key = t.key and input_ids.path = t.path and input_ids.locale = t.locale
`

func (q *Queries) FindTranslations(ctx context.Context, keys pgtype.JSONB) ([]Translation, error) {
	rows, err := q.db.Query(ctx, findTranslations, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Translation
	for rows.Next() {
		var i Translation
		if err := rows.Scan(
			&i.Key,
			&i.Path,
			&i.Locale,
			&i.Value,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PseudoID,
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

const saveTranslation = `-- name: SaveTranslation :one
insert into "translation"."translation"
    (key, path, value, locale)
values ($1, $2, $3, $4)
ON CONFLICT (key, path, locale) DO UPDATE
    SET value      = excluded.value,
        updated_at = now()
RETURNING key, path, locale, value, created_at, updated_at, pseudo_id
`

type SaveTranslationParams struct {
	Key    string `db:"key" json:"key"`
	Path   Path   `db:"path" json:"path"`
	Value  string `db:"value" json:"value"`
	Locale Locale `db:"locale" json:"locale"`
}

func (q *Queries) SaveTranslation(ctx context.Context, arg SaveTranslationParams) (Translation, error) {
	row := q.db.QueryRow(ctx, saveTranslation,
		arg.Key,
		arg.Path,
		arg.Value,
		arg.Locale,
	)
	var i Translation
	err := row.Scan(
		&i.Key,
		&i.Path,
		&i.Locale,
		&i.Value,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PseudoID,
	)
	return i, err
}
