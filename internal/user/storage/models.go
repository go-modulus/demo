// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package storage

import (
	"time"

	uuid "github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
)

type User struct {
	ID           uuid.UUID    `db:"id" json:"id"`
	Name         string       `db:"name" json:"name"`
	Email        string       `db:"email" json:"email"`
	RegisteredAt time.Time    `db:"registered_at" json:"registeredAt"`
	Settings     pgtype.JSONB `db:"settings" json:"settings"`
	Contacts     []string     `db:"contacts" json:"contacts"`
}
