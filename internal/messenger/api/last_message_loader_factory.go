package api

import (
	"context"
	"demo/internal/messenger/domain"
	"github.com/gofrs/uuid"
	"github.com/graph-gophers/dataloader/v7"
)

type LastMessageLoaderFactory struct {
	messages domain.MessageRepository
}

func NewLastMessageLoaderFactory(messages domain.MessageRepository) *LastMessageLoaderFactory {
	return &LastMessageLoaderFactory{messages: messages}
}

func (l *LastMessageLoaderFactory) Create() *dataloader.Loader[uuid.UUID, *domain.Message] {
	return dataloader.NewBatchedLoader(
		func(ctx context.Context, conversationIds []uuid.UUID) []*dataloader.Result[*domain.Message] {
			messages, err := l.messages.FindLastMessages(ctx, conversationIds)
			if err != nil {
				return []*dataloader.Result[*domain.Message]{
					{Error: err},
				}
			}

			result := make([]*dataloader.Result[*domain.Message], len(conversationIds))
			for i, conversationId := range conversationIds {
				result[i] = &dataloader.Result[*domain.Message]{
					Data: messages[conversationId],
				}
			}
			return result
		},
	)
}
