// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package persistence

import (
	"database/sql"
	"time"

	"demo/internal/messenger/domain"
	"github.com/gofrs/uuid"
)

type MessengerConversation struct {
	ID         uuid.UUID
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

type MessengerDraft struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	AuthorID       uuid.UUID
	Text           sql.NullString
	TextParts      domain.RichTextParts
}

type MessengerMessage struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Text           sql.NullString
	TextParts      domain.RichTextParts
	Status         string
	Type           string
	UpdatedAt      time.Time
	CreatedAt      time.Time
}