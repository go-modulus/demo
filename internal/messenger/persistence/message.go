package persistence

import (
	"context"
	"database/sql"
	"demo/internal/messenger/domain"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid"
	"time"
)

type MessageRepository struct {
	db      *sql.DB
	queries *Queries
}

func NewMessageRepository(db *sql.DB, queries *Queries) *MessageRepository {
	return &MessageRepository{db: db, queries: queries}
}

func (m *MessageRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	message, err := m.queries.GetMessage(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrMessageNotFound{MessageId: id}
	}
	if err != nil {
		return nil, err
	}

	return m.mapToDomain(message), nil
}

func (m *MessageRepository) FindLastMessages(ctx context.Context, conversationIds []uuid.UUID) (map[uuid.UUID]*domain.Message, error) {
	messages, err := m.queries.FindLastMessages(ctx, conversationIds)
	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]*domain.Message)
	for _, message := range messages {
		result[message.ConversationID] = m.mapToDomain(message)
	}

	return result, nil
}

func (m *MessageRepository) Paginate(
	ctx context.Context,
	conversationId uuid.UUID,
	first int,
	after sql.NullString,
) (*domain.Page[*domain.Message], error) {
	params := PaginateMessagesParams{
		ConversationID: conversationId,
		First:          int32(first) + 1,
	}

	if after.Valid {
		cursor, err := m.decodeCursor(after.String)
		if err != nil {
			return nil, err
		}

		if rawCreatedAt, ok := cursor["createdAt"]; ok {
			params.AfterCreatedAt = sql.NullTime{
				Time:  time.UnixMicro(int64(rawCreatedAt.(float64))),
				Valid: true,
			}

			if rawId, ok := cursor["id"]; ok {
				id, err := uuid.FromString(rawId.(string))
				if err != nil {
					return nil, err
				}

				params.AfterID = uuid.NullUUID{
					UUID:  id,
					Valid: true,
				}
			}
		}
	}

	messages, err := m.queries.PaginateMessages(ctx, params)
	if err != nil {
		return nil, err
	}

	edges := make([]domain.Edge[*domain.Message], 0, len(messages))
	for i, message := range messages {
		if i == first {
			break
		}

		cursor, err := m.encodeCursor(map[string]any{
			"createdAt": message.CreatedAt.UnixMicro(),
			"id":        message.ID.String(),
		})
		if err != nil {
			return nil, err
		}

		edges = append(edges, domain.Edge[*domain.Message]{
			Cursor: cursor,
			Node:   m.mapToDomain(message),
		})
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		if len(edges) > 0 {
			startCursor = &edges[0].Cursor
			endCursor = &edges[len(edges)-1].Cursor
		}
	}

	return &domain.Page[*domain.Message]{
		Edges:           edges,
		HasNextPage:     len(messages) > len(edges),
		HasPreviousPage: after.Valid,
		StartCursor:     startCursor,
		EndCursor:       endCursor,
	}, nil
}

func (m *MessageRepository) encodeCursor(cursor map[string]any) (string, error) {
	rawJson, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(rawJson), nil
}

func (m *MessageRepository) decodeCursor(rawCursor string) (map[string]any, error) {
	rawJson, err := base64.StdEncoding.DecodeString(rawCursor)
	if err != nil {
		return nil, err
	}

	cursor := make(map[string]any)
	err = json.Unmarshal(rawJson, &cursor)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func (m *MessageRepository) Update(ctx context.Context, id uuid.UUID, updateFunc func(ctx context.Context, message *domain.Message) (*domain.Message, error)) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := m.queries.WithTx(tx)
	rawMessage, err := queries.FindMessageForUpdate(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrMessageNotFound{MessageId: id}
	}
	if err != nil {
		return err
	}

	message := m.mapToDomain(rawMessage)
	message, err = updateFunc(ctx, message)
	if err != nil {
		return err
	}

	err = queries.UpdateMessage(ctx, UpdateMessageParams{
		ID:        message.Id,
		Text:      message.Text.Text,
		TextParts: message.Text.Parts,
		UpdatedAt: message.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *MessageRepository) Add(ctx context.Context, message *domain.Message) error {
	return m.queries.CreateMessage(ctx, CreateMessageParams{
		ID:             message.Id,
		ConversationID: message.ConversationId,
		SenderID:       message.SenderId,
		Text:           message.Text.Text,
		TextParts:      message.Text.Parts,
		Status:         "created",
		Type:           "text",
		UpdatedAt:      message.UpdatedAt,
		CreatedAt:      message.CreatedAt,
	})
}

func (m *MessageRepository) mapToDomain(message MessengerMessage) *domain.Message {
	return &domain.Message{
		Id:             message.ID,
		ConversationId: message.ConversationID,
		SenderId:       message.SenderID,
		Text: domain.RichText{
			Text:  message.Text,
			Parts: message.TextParts,
		},
		UpdatedAt: message.UpdatedAt,
		CreatedAt: message.CreatedAt,
	}
}
