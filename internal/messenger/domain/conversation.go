package domain

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

type ErrConversationNotFound struct {
	ConversationId uuid.UUID
}

func (e ErrConversationNotFound) Error() string {
	return fmt.Sprintf("conversation with id %s not found", e.ConversationId)
}

type Conversation struct {
	Id         uuid.UUID
	SenderId   uuid.UUID
	ReceiverId uuid.UUID
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

func (c Conversation) IsParticipant(personId uuid.UUID) bool {
	return c.SenderId == personId || c.ReceiverId == personId
}

type ConversationRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*Conversation, error)
	Paginate(ctx context.Context, viewer uuid.UUID, first int, after sql.NullString) (*Page[*Conversation], error)

	GetOrCreate(ctx context.Context, sender uuid.UUID, receiver uuid.UUID) (*Conversation, error)
}

type ConversationCreator struct {
	conversations ConversationRepository
}

func NewConversationCreator(conversations ConversationRepository) *ConversationCreator {
	return &ConversationCreator{conversations: conversations}
}

func (c *ConversationCreator) Create(ctx context.Context, sender uuid.UUID, receiver uuid.UUID) (*Conversation, error) {
	conversation, err := c.conversations.GetOrCreate(ctx, sender, receiver)
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

type ConversationViewer struct {
	conversations ConversationRepository
	policy        *ConversationPolicy
}

func NewConversationViewer(conversations ConversationRepository, policy *ConversationPolicy) *ConversationViewer {
	return &ConversationViewer{conversations: conversations, policy: policy}
}

func (c *ConversationViewer) View(ctx context.Context, viewer uuid.UUID, id uuid.UUID) (*Conversation, error) {
	conversation, err := c.conversations.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !c.policy.CanView(viewer, conversation) {
		return nil, ErrConversationNotFound{ConversationId: id}
	}

	return conversation, nil
}

func (c *ConversationViewer) Paginate(
	ctx context.Context,
	viewer uuid.UUID,
	first int,
	after sql.NullString,
) (*Page[*Conversation], error) {
	page, err := c.conversations.Paginate(ctx, viewer, first, after)
	if err != nil {
		return nil, err
	}

	return page, nil
}

type ConversationPolicy struct{}

func NewConversationPolicy() *ConversationPolicy {
	return &ConversationPolicy{}
}

func (c *ConversationPolicy) CanView(personId uuid.UUID, conversation *Conversation) bool {
	if conversation == nil {
		return false
	}

	return conversation.IsParticipant(personId)
}

func (c *ConversationPolicy) CanSendMessages(personId uuid.UUID, conversation *Conversation) bool {
	if conversation == nil {
		return false
	}

	return conversation.IsParticipant(personId)
}
