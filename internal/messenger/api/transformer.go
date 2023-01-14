package api

import (
	"demo/graph/model"
	"demo/internal/messenger/domain"
	"fmt"
	"github.com/gofrs/uuid"
)

type Transformer struct {
}

func NewTransformer() *Transformer {
	return &Transformer{}
}

func (t *Transformer) TransformConversationPage(page *domain.Page[*domain.Conversation]) *model.ConversationList {
	edges := make([]*model.ConversationEdge, len(page.Edges))
	for i, edge := range page.Edges {
		edges[i] = &model.ConversationEdge{
			Cursor: edge.Cursor,
			Node:   t.TransformConversation(edge.Node),
		}
	}

	startCursor, endCursor := "", ""
	if page.StartCursor != nil {
		startCursor = *page.StartCursor
	}
	if page.EndCursor != nil {
		endCursor = *page.EndCursor
	}

	return &model.ConversationList{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     page.HasNextPage,
			HasPreviousPage: page.HasPreviousPage,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
	}
}

func (t *Transformer) TransformMessagePage(page *domain.Page[*domain.Message]) *model.MessageList {
	edges := make([]*model.MessageEdge, len(page.Edges))
	for i, edge := range page.Edges {
		edges[i] = &model.MessageEdge{
			Cursor: edge.Cursor,
			Node:   t.TransformMessage(edge.Node),
		}
	}

	startCursor, endCursor := "", ""
	if page.StartCursor != nil {
		startCursor = *page.StartCursor
	}
	if page.EndCursor != nil {
		endCursor = *page.EndCursor
	}

	return &model.MessageList{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     page.HasNextPage,
			HasPreviousPage: page.HasPreviousPage,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
	}
}

func (t *Transformer) TransformConversation(conversation *domain.Conversation) model.Conversation {
	return t.TransformOneToOneConversation(conversation)
}

func (t *Transformer) TransformOneToOneConversation(conversation *domain.Conversation) *model.OneToOneConversation {
	return &model.OneToOneConversation{
		ID:        conversation.Id,
		CreatedAt: conversation.CreatedAt,
	}
}

func (t *Transformer) TransformDraft(draft *domain.Draft) *model.Draft {
	var messageId *uuid.UUID
	if draft.MessageId.Valid {
		messageId = &draft.MessageId.UUID
	}

	return &model.Draft{
		ConversationID: draft.ConversationId,
		MessageID:      messageId,
		RichText:       t.transformRichText(draft.Text),
	}
}

func (t *Transformer) TransformMessage(message *domain.Message) model.Message {
	return t.TransformTextMessage(message)
}

func (t *Transformer) TransformTextMessage(message *domain.Message) *model.TextMessage {
	return &model.TextMessage{
		ID:             message.Id,
		ConversationID: message.ConversationId,
		RichText:       t.transformRichText(message.Text),
		UpdatedAt:      message.UpdatedAt,
		CreatedAt:      message.CreatedAt,
	}
}

func (t *Transformer) transformRichText(text domain.RichText) *model.RichText {
	if !text.Text.Valid {
		return nil
	}

	richText := &model.RichText{
		Text:  &text.Text.String,
		Parts: make([]model.RichTextPart, len(text.Parts)),
	}

	for i, part := range text.Parts {
		switch part := part.(type) {
		case domain.PlainRichText:
			richText.Parts[i] = &model.PlainRichText{
				Text: part.Text(),
			}
		default:
			panic(fmt.Sprintf("unknown rich text part type: %s", part.Type()))
		}
	}

	return richText
}
