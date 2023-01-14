package domain

import (
	"context"
	"github.com/gofrs/uuid"
)

type Draft struct {
	Id             uuid.UUID
	ConversationId uuid.UUID
	AuthorId       uuid.UUID
	MessageId      uuid.NullUUID
	Text           RichText
}

type DraftRepository interface {
	FindDrafts(
		ctx context.Context,
		authorId uuid.UUID,
		conversationIds []uuid.UUID,
	) (map[uuid.UUID]*Draft, error)
	FindOrCreateByConversationAndAuthor(
		ctx context.Context,
		conversationId uuid.UUID,
		authorId uuid.UUID,
	) (*Draft, error)
	Update(
		ctx context.Context,
		id uuid.UUID,
		updateFunc func(ctx context.Context, draft *Draft) (*Draft, error),
	) error
	RemoveByConversationAndAuthor(
		ctx context.Context,
		conversationId uuid.UUID,
		authorId uuid.UUID,
	) error
}

type SaveDraft struct {
	ConversationId uuid.UUID
	SenderId       uuid.UUID
	MessageId      uuid.NullUUID
	Text           string
}

type DraftSaver struct {
	conversations      ConversationRepository
	conversationPolicy *ConversationPolicy
	messages           MessageRepository
	messagePolicy      *MessagePolicy
	richTextParser     *RichTextParser
	drafts             DraftRepository
}

func NewDraftSaver(
	conversations ConversationRepository,
	conversationPolicy *ConversationPolicy,
	messages MessageRepository,
	messagePolicy *MessagePolicy,
	richTextParser *RichTextParser,
	drafts DraftRepository,
) *DraftSaver {
	return &DraftSaver{
		conversations:      conversations,
		conversationPolicy: conversationPolicy,
		messages:           messages,
		messagePolicy:      messagePolicy,
		richTextParser:     richTextParser,
		drafts:             drafts,
	}
}

func (s *DraftSaver) Save(ctx context.Context, dto SaveDraft) (*Draft, error) {
	conversation, err := s.conversations.Get(ctx, dto.ConversationId)
	if err != nil {
		return nil, err
	}
	if !s.conversationPolicy.CanView(dto.SenderId, conversation) {
		return nil, ErrConversationNotFound{ConversationId: dto.ConversationId}
	}

	var message *Message
	if dto.MessageId.Valid {
		message, err = s.messages.Get(ctx, dto.MessageId.UUID)
		if err != nil {
			return nil, err
		}
		if !s.messagePolicy.CanEdit(dto.SenderId, message) {
			return nil, ErrMessageNotFound{MessageId: dto.MessageId.UUID}
		}
		if message.ConversationId != conversation.Id {
			return nil, ErrMessageNotFound{MessageId: dto.MessageId.UUID}
		}
	}

	richText, err := s.richTextParser.Parse(dto.Text)
	if err != nil {
		return nil, err
	}

	draftToReturn, err := s.drafts.FindOrCreateByConversationAndAuthor(
		ctx,
		conversation.Id,
		dto.SenderId,
	)
	if err != nil {
		return nil, err
	}

	err = s.drafts.Update(
		ctx,
		draftToReturn.Id,
		func(ctx context.Context, draft *Draft) (*Draft, error) {
			if message != nil {
				draft.MessageId = uuid.NullUUID{UUID: message.Id, Valid: true}
			} else {
				draft.MessageId = uuid.NullUUID{}
			}

			draft.Text = richText
			draftToReturn = draft

			return draft, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return draftToReturn, nil
}

type RemoveDraft struct {
	ConversationId uuid.UUID
	AuthorId       uuid.UUID
}

type DraftRemover struct {
	drafts DraftRepository
}

func NewDraftRemover(drafts DraftRepository) *DraftRemover {
	return &DraftRemover{drafts: drafts}
}

func (s *DraftRemover) Remove(ctx context.Context, dto RemoveDraft) error {
	return s.drafts.RemoveByConversationAndAuthor(ctx, dto.ConversationId, dto.AuthorId)
}
