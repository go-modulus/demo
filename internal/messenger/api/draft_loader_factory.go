package api

import (
	"context"
	"demo/internal/auth"
	"demo/internal/messenger/domain"
	"github.com/gofrs/uuid"
	"github.com/graph-gophers/dataloader/v7"
)

type DraftLoaderFactory struct {
	drafts domain.DraftRepository
}

func NewDraftLoaderFactory(drafts domain.DraftRepository) *DraftLoaderFactory {
	return &DraftLoaderFactory{drafts: drafts}
}

func (l *DraftLoaderFactory) Create() *dataloader.Loader[uuid.UUID, *domain.Draft] {
	return dataloader.NewBatchedLoader(
		func(ctx context.Context, conversationIds []uuid.UUID) []*dataloader.Result[*domain.Draft] {
			performer := auth.PerformerFromContext(ctx)
			if !performer.Valid {
				return nil
			}

			drafts, err := l.drafts.FindDrafts(ctx, performer.Value.Id, conversationIds)
			if err != nil {
				return []*dataloader.Result[*domain.Draft]{
					{Error: err},
				}
			}

			result := make([]*dataloader.Result[*domain.Draft], len(conversationIds))
			for i, conversationId := range conversationIds {
				result[i] = &dataloader.Result[*domain.Draft]{
					Data: drafts[conversationId],
				}
			}
			return result
		},
	)
}
