package domain

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

type ErrMessageNotFound struct {
	MessageId uuid.UUID
}

func (e ErrMessageNotFound) Error() string {
	return fmt.Sprintf("message with id %s not found", e.MessageId)
}

var (
	MessageStatusUnknown = MessageStatus{""}
	MessageStatusNew     = MessageStatus{"confirmed"}
	MessageStatusDeleted = MessageStatus{"deleted"}
)

type MessageStatus struct {
	value string
}

func NewMessageStatus(value string) MessageStatus {
	switch value {
	case MessageStatusNew.value:
		return MessageStatusNew
	case MessageStatusDeleted.value:
		return MessageStatusDeleted
	}

	return MessageStatusUnknown
}

type MessageCreated struct {
	MessageId      uuid.UUID
	ConversationId uuid.UUID
	SenderId       uuid.UUID
	Text           RichText
	CreatedAt      time.Time
	BaseEvent
}

type Message struct {
	Id             uuid.UUID
	ConversationId uuid.UUID
	SenderId       uuid.UUID
	Text           RichText
	Status         MessageStatus
	UpdatedAt      time.Time
	CreatedAt      time.Time
	EventCollector
}

func (m *Message) MarkAsDeleted() {
	if m.Status == MessageStatusDeleted {
		return
	}

	m.Status = MessageStatusDeleted
	m.UpdatedAt = time.Now()
}

func (m *Message) IsSender(personId uuid.UUID) bool {
	return m.SenderId == personId
}

func NewMessage(conversationID uuid.UUID, senderID uuid.UUID, richText RichText) *Message {
	now := time.Now()

	message := &Message{
		Id:             uuid.Must(uuid.NewV6()),
		ConversationId: conversationID,
		SenderId:       senderID,
		Text:           richText,
		UpdatedAt:      now,
		CreatedAt:      now,
	}

	message.recordThat(&MessageCreated{
		MessageId:      message.Id,
		ConversationId: message.ConversationId,
		SenderId:       message.SenderId,
		Text:           message.Text,
		CreatedAt:      message.CreatedAt,
	})

	return message
}

type MessageRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*Message, error)
	FindLastMessages(ctx context.Context, conversationIds []uuid.UUID) (map[uuid.UUID]*Message, error)
	Paginate(ctx context.Context, conversationId uuid.UUID, first int, after sql.NullString) (*Page[*Message], error)
	Update(
		ctx context.Context,
		id uuid.UUID,
		updateFunc func(ctx context.Context, draft *Message) (*Message, error),
	) error
	Add(ctx context.Context, message *Message) error
}

type CreateMessage struct {
	PerformerId    uuid.UUID
	ConversationId uuid.UUID
	Text           string
}

type MessageCreator struct {
	richTextParser *RichTextParser
	messages       MessageRepository
}

func NewMessageCreator(richTextParser *RichTextParser, messages MessageRepository) *MessageCreator {
	return &MessageCreator{richTextParser: richTextParser, messages: messages}
}

func (c *MessageCreator) Create(ctx context.Context, dto CreateMessage) (*Message, error) {
	richText, err := c.richTextParser.Parse(dto.Text)
	if err != nil {
		return nil, err
	}

	message := NewMessage(dto.ConversationId, dto.PerformerId, richText)

	if err := c.messages.Add(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

type EditMessage struct {
	PerformerId uuid.UUID
	MessageId   uuid.UUID
	Text        string
}

type MessageEditor struct {
	richTextParser *RichTextParser
	messages       MessageRepository
}

func NewMessageEditor(richTextParser *RichTextParser, messages MessageRepository) *MessageEditor {
	return &MessageEditor{richTextParser: richTextParser, messages: messages}
}

func (c *MessageEditor) Edit(ctx context.Context, dto EditMessage) (*Message, error) {
	richText, err := c.richTextParser.Parse(dto.Text)
	if err != nil {
		return nil, err
	}

	var messageToReturn *Message
	err = c.messages.Update(
		ctx,
		dto.MessageId,
		func(ctx context.Context, message *Message) (*Message, error) {
			message.Text = richText
			message.UpdatedAt = time.Now()
			messageToReturn = message
			return message, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return messageToReturn, nil
}

type RemoveMessage struct {
	PerformerId uuid.UUID
	MessageId   uuid.UUID
}

type MessageRemover struct {
	messages MessageRepository
	policy   *MessagePolicy
}

func NewMessageRemover(messages MessageRepository, policy *MessagePolicy) *MessageRemover {
	return &MessageRemover{messages: messages, policy: policy}
}

func (r *MessageRemover) Remove(ctx context.Context, dto RemoveMessage) error {
	return r.messages.Update(
		ctx,
		dto.MessageId,
		func(ctx context.Context, message *Message) (*Message, error) {
			if !r.policy.CanDelete(dto.PerformerId, message) {
				return nil, ErrMessageNotFound{MessageId: dto.MessageId}
			}

			message.MarkAsDeleted()

			return message, nil
		},
	)
}

type MessageViewer struct {
	conversations      ConversationRepository
	conversationPolicy *ConversationPolicy
	messages           MessageRepository
}

func NewMessageViewer(
	conversations ConversationRepository,
	conversationPolicy *ConversationPolicy,
	messages MessageRepository,
) *MessageViewer {
	return &MessageViewer{
		conversations:      conversations,
		conversationPolicy: conversationPolicy,
		messages:           messages,
	}
}

func (c *MessageViewer) Paginate(
	ctx context.Context,
	performerId uuid.UUID,
	conversationId uuid.UUID,
	first int,
	after sql.NullString,
) (*Page[*Message], error) {
	conversation, err := c.conversations.Get(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	if !c.conversationPolicy.CanView(performerId, conversation) {
		return nil, ErrConversationNotFound{ConversationId: conversationId}
	}

	page, err := c.messages.Paginate(ctx, conversationId, first, after)
	if err != nil {
		return nil, err
	}

	return page, nil
}

type MessagePolicy struct{}

func NewMessagePolicy() *MessagePolicy {
	return &MessagePolicy{}
}

func (c *MessagePolicy) CanEdit(personId uuid.UUID, message *Message) bool {
	if message == nil {
		return false
	}

	return message.IsSender(personId)
}

func (c *MessagePolicy) CanDelete(personId uuid.UUID, message *Message) bool {
	if message == nil {
		return false
	}

	return message.IsSender(personId)
}
