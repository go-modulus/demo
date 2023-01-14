package api

import (
	"context"
	"database/sql"
	"demo/graph/model"
	"demo/internal/auth"
	"demo/internal/errors"
	"demo/internal/graphql"
	"demo/internal/messenger/domain"
	"fmt"
	"github.com/gofrs/uuid"
)

type Resolver struct {
	conversationViewer       *domain.ConversationViewer
	conversationCreator      *domain.ConversationCreator
	draftSaver               *domain.DraftSaver
	draftRemover             *domain.DraftRemover
	messageViewer            *domain.MessageViewer
	messageCreator           *domain.MessageCreator
	messageEditor            *domain.MessageEditor
	messageRemover           *domain.MessageRemover
	draftLoaderFactory       *DraftLoaderFactory
	lastMessageLoaderFactory *LastMessageLoaderFactory
	transformer              *Transformer
}

func NewResolver(
	conversationViewer *domain.ConversationViewer,
	conversationCreator *domain.ConversationCreator,
	draftSaver *domain.DraftSaver,
	draftRemover *domain.DraftRemover,
	messageViewer *domain.MessageViewer,
	messageCreator *domain.MessageCreator,
	messageEditor *domain.MessageEditor,
	messageRemover *domain.MessageRemover,
	draftLoaderFactory *DraftLoaderFactory,
	lastMessageLoaderFactory *LastMessageLoaderFactory,
	transformer *Transformer,
) *Resolver {
	return &Resolver{
		conversationViewer:       conversationViewer,
		conversationCreator:      conversationCreator,
		draftSaver:               draftSaver,
		draftRemover:             draftRemover,
		messageViewer:            messageViewer,
		messageCreator:           messageCreator,
		messageEditor:            messageEditor,
		messageRemover:           messageRemover,
		draftLoaderFactory:       draftLoaderFactory,
		lastMessageLoaderFactory: lastMessageLoaderFactory,
		transformer:              transformer,
	}
}

func (r *Resolver) CreateOneToOneConversation(ctx context.Context, receiverID uuid.UUID) (model.CreateOneToOneConversationResult, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return model.NewErrUnauthorized(), nil
	}

	conversation, err := r.conversationCreator.Create(ctx, performer.Value.Id, receiverID)
	if err != nil {
		return nil, err
	}

	return r.transformer.TransformOneToOneConversation(conversation), nil
}

func (r *Resolver) GetConversation(ctx context.Context, id uuid.UUID) (model.ConversationResult, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return model.NewErrUnauthorized(), nil
	}

	conversation, err := r.conversationViewer.View(ctx, performer.Value.Id, id)
	if err != nil {
		var errConversationNotFound domain.ErrConversationNotFound
		if errors.As(err, &errConversationNotFound) {
			return model.ErrUnknownConversation{
				Message: fmt.Sprintf(
					"Conversation#%s not found",
					errConversationNotFound.ConversationId,
				),
			}, nil
		}

		return nil, err
	}

	return &model.ConversationBox{
		Value: r.transformer.TransformConversation(conversation),
	}, nil
}

func (r *Resolver) MyConversations(ctx context.Context, first int, rawAfter *string) (model.MyConversationsResult, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return model.NewErrUnauthorized(), nil
	}

	var after sql.NullString
	if rawAfter != nil {
		after = sql.NullString{
			String: *rawAfter,
			Valid:  true,
		}
	}

	page, err := r.conversationViewer.Paginate(ctx, performer.Value.Id, first, after)
	if err != nil {
		return nil, err
	}

	return r.transformer.TransformConversationPage(page), nil
}

func (r *Resolver) Draft(ctx context.Context, conversation *model.OneToOneConversation) (*model.Draft, error) {
	loader := graphql.GetLoader[uuid.UUID, *domain.Draft](ctx, r.draftLoaderFactory)
	draft, err := loader.Load(ctx, conversation.ID)()
	if err != nil || draft == nil {
		return nil, err
	}
	return r.transformer.TransformDraft(draft), nil
}

func (r *Resolver) LastMessage(ctx context.Context, conversation *model.OneToOneConversation) (*model.TextMessage, error) {
	loader := graphql.GetLoader[uuid.UUID, *domain.Message](ctx, r.lastMessageLoaderFactory)
	message, err := loader.Load(ctx, conversation.ID)()
	if err != nil || message == nil {
		return nil, err
	}
	return r.transformer.TransformTextMessage(message), nil
}

func (r *Resolver) SaveDraft(
	ctx context.Context,
	conversationId uuid.UUID,
	rawMessageId *uuid.UUID,
	rawText *string,
) (model.SaveDraftResult, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return model.NewErrUnauthorized(), nil
	}

	var messageId uuid.NullUUID
	if rawMessageId != nil {
		messageId = uuid.NullUUID{
			UUID:  *rawMessageId,
			Valid: true,
		}
	}

	var text string
	if rawText != nil {
		text = *rawText
	}

	message, err := r.draftSaver.Save(ctx, domain.SaveDraft{
		ConversationId: conversationId,
		SenderId:       performer.Value.Id,
		MessageId:      messageId,
		Text:           text,
	})
	if err != nil {
		var errConversationNotFound domain.ErrConversationNotFound
		if errors.As(err, &errConversationNotFound) {
			return model.ErrUnknownConversation{
				Message: fmt.Sprintf(
					"Conversation#%s not found",
					errConversationNotFound.ConversationId,
				),
			}, nil
		}

		var errMessageNotFound domain.ErrMessageNotFound
		if errors.As(err, &errMessageNotFound) {
			return model.ErrUnknownMessage{
				Message: fmt.Sprintf(
					"Message#%s not found",
					errMessageNotFound.MessageId,
				),
			}, nil
		}

		return nil, err
	}

	return r.transformer.TransformDraft(message), nil
}

func (r *Resolver) RemoveDraft(ctx context.Context, conversationId uuid.UUID) (model.RemoveDraftResult, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return model.NewErrUnauthorized(), nil
	}

	err := r.draftRemover.Remove(ctx, domain.RemoveDraft{
		ConversationId: conversationId,
		AuthorId:       performer.Value.Id,
	})

	return nil, err
}

func (r *Resolver) Messages(
	ctx context.Context,
	conversationId uuid.UUID,
	first int,
	rawAfter *string,
) (model.MessagesResult, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return model.NewErrUnauthorized(), nil
	}

	if conversationId.IsNil() {
		return model.ErrInvalidInput{
			Message: "Invalid input",
			Fields: []*model.InvalidField{
				{
					Field: "conversationId",
					Rule: model.UUIDRule{
						Message: "Invalid UUID",
					},
				},
			},
		}, nil
	}

	var after sql.NullString
	if rawAfter != nil {
		after = sql.NullString{
			String: *rawAfter,
			Valid:  true,
		}
	}

	page, err := r.messageViewer.Paginate(
		ctx,
		performer.Value.Id,
		conversationId,
		first,
		after,
	)
	if err != nil {
		var errConversationNotFound domain.ErrConversationNotFound
		if errors.As(err, &errConversationNotFound) {
			return model.ErrUnknownConversation{
				Message: fmt.Sprintf(
					"Conversation#%s not found",
					errConversationNotFound.ConversationId,
				),
			}, nil
		}

		return nil, err
	}

	return r.transformer.TransformMessagePage(page), nil
}

func (r *Resolver) CreateMessage(
	ctx context.Context,
	conversationId uuid.UUID,
	rawText *string,
) (model.CreateMessageResult, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return model.NewErrUnauthorized(), nil
	}

	var text string
	if rawText != nil {
		text = *rawText
	}

	message, err := r.messageCreator.Create(ctx, domain.CreateMessage{
		PerformerId:    performer.Value.Id,
		ConversationId: conversationId,
		Text:           text,
	})
	if err != nil {
		return nil, err
	}

	return r.transformer.TransformTextMessage(message), nil
}

func (r *Resolver) EditMessage(
	ctx context.Context,
	messageId uuid.UUID,
	rawText *string,
) (model.EditMessageResult, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return model.NewErrUnauthorized(), nil
	}

	var text string
	if rawText != nil {
		text = *rawText
	}

	message, err := r.messageEditor.Edit(ctx, domain.EditMessage{
		PerformerId: performer.Value.Id,
		MessageId:   messageId,
		Text:        text,
	})
	if err != nil {
		var errMessageNotFound domain.ErrMessageNotFound
		if errors.As(err, &errMessageNotFound) {
			return model.ErrUnknownMessage{
				Message: fmt.Sprintf(
					"Message#%s not found",
					errMessageNotFound.MessageId,
				),
			}, nil
		}

		return nil, err
	}

	return r.transformer.TransformTextMessage(message), nil
}

func (r *Resolver) DeleteMessage(
	ctx context.Context,
	messageId uuid.UUID,
) (model.DeleteMessageResult, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return model.NewErrUnauthorized(), nil
	}

	err := r.messageRemover.Remove(ctx, domain.RemoveMessage{
		PerformerId: performer.Value.Id,
		MessageId:   messageId,
	})
	if err != nil {
		var errMessageNotFound domain.ErrMessageNotFound
		if errors.As(err, &errMessageNotFound) {
			return model.ErrUnknownMessage{
				Message: fmt.Sprintf(
					"Message#%s not found",
					errMessageNotFound.MessageId,
				),
			}, nil
		}

		return nil, err
	}

	return model.VoidBox{}, nil
}
