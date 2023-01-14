package persistence

import (
	"context"
	"database/sql"
	"demo/internal/errors"
	"demo/internal/messenger/domain"
	"encoding/base64"
	"encoding/json"
	"github.com/gofrs/uuid"
	"time"
)

var _ domain.ConversationRepository = &ConversationRepository{}

type ConversationRepository struct {
	queries *Queries
}

func NewConversationRepository(queries *Queries) *ConversationRepository {
	return &ConversationRepository{queries: queries}
}

func (c *ConversationRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
	conversation, err := c.queries.GetConversation(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrConversationNotFound{ConversationId: id}
	}
	if err != nil {
		return nil, err
	}

	return &domain.Conversation{
		Id:         conversation.ID,
		SenderId:   conversation.SenderID,
		ReceiverId: conversation.ReceiverID,
		UpdatedAt:  conversation.UpdatedAt,
		CreatedAt:  conversation.CreatedAt,
	}, nil
}

func (c *ConversationRepository) Paginate(
	ctx context.Context,
	viewer uuid.UUID,
	first int,
	after sql.NullString,
) (*domain.Page[*domain.Conversation], error) {
	params := PaginateMyConversationsParams{
		ViewerID: viewer,
		First:    int32(first) + 1,
	}

	if after.Valid {
		cursor, err := c.decodeCursor(after.String)
		if err != nil {
			return nil, err
		}

		if rawUpdatedAt, ok := cursor["updatedAt"]; ok {
			params.AfterUpdatedAt = sql.NullTime{
				Time:  time.UnixMicro(int64(rawUpdatedAt.(float64))),
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

	conversations, err := c.queries.PaginateMyConversations(ctx, params)
	if err != nil {
		return nil, err
	}

	edges := make([]domain.Edge[*domain.Conversation], 0, len(conversations))
	for i, conversation := range conversations {
		if i == first {
			break
		}

		cursor, err := c.encodeCursor(map[string]any{
			"updatedAt": conversation.UpdatedAt.UnixMicro(),
			"id":        conversation.ID.String(),
		})
		if err != nil {
			return nil, err
		}

		edges = append(edges, domain.Edge[*domain.Conversation]{
			Cursor: cursor,
			Node:   c.mapToDomain(conversation),
		})
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		if len(edges) > 0 {
			startCursor = &edges[0].Cursor
			endCursor = &edges[len(edges)-1].Cursor
		}
	}

	return &domain.Page[*domain.Conversation]{
		Edges:           edges,
		HasNextPage:     len(conversations) > len(edges),
		HasPreviousPage: after.Valid,
		StartCursor:     startCursor,
		EndCursor:       endCursor,
	}, nil
}

func (c *ConversationRepository) GetOrCreate(ctx context.Context, sender uuid.UUID, receiver uuid.UUID) (*domain.Conversation, error) {
	conversation, err := c.queries.CreateOrGetConversation(ctx, CreateOrGetConversationParams{
		ID:         uuid.Must(uuid.NewV6()),
		SenderID:   sender,
		ReceiverID: receiver,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return &domain.Conversation{
		Id:         conversation.ID,
		SenderId:   conversation.SenderID,
		ReceiverId: conversation.ReceiverID,
		UpdatedAt:  conversation.UpdatedAt,
		CreatedAt:  conversation.CreatedAt,
	}, nil
}

func (c *ConversationRepository) encodeCursor(cursor map[string]any) (string, error) {
	rawJson, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(rawJson), nil
}

func (c *ConversationRepository) decodeCursor(rawCursor string) (map[string]any, error) {
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

func (c *ConversationRepository) mapToDomain(conversation MessengerConversation) *domain.Conversation {
	return &domain.Conversation{
		Id:         conversation.ID,
		SenderId:   conversation.SenderID,
		ReceiverId: conversation.ReceiverID,
		UpdatedAt:  conversation.UpdatedAt,
		CreatedAt:  conversation.CreatedAt,
	}
}
