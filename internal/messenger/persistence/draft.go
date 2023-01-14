package persistence

import (
	"context"
	"database/sql"
	"demo/internal/messenger/domain"
	"github.com/gofrs/uuid"
)

type DraftRepository struct {
	db      *sql.DB
	queries *Queries
}

func NewDraftRepository(db *sql.DB, queries *Queries) *DraftRepository {
	return &DraftRepository{db: db, queries: queries}
}

func (d *DraftRepository) FindDrafts(
	ctx context.Context,
	authorId uuid.UUID,
	conversationIds []uuid.UUID,
) (map[uuid.UUID]*domain.Draft, error) {
	drafts, err := d.queries.FindDrafts(ctx, FindDraftsParams{
		AuthorID:        authorId,
		ConversationIds: conversationIds,
	})
	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]*domain.Draft)
	for _, draft := range drafts {
		result[draft.ConversationID] = d.mapToDomain(draft)
	}

	return result, nil
}

func (d *DraftRepository) FindOrCreateByConversationAndAuthor(ctx context.Context, conversationId uuid.UUID, authorId uuid.UUID) (*domain.Draft, error) {
	draft, err := d.queries.FindOrCreateDraft(ctx, FindOrCreateDraftParams{
		ID:             uuid.Must(uuid.NewV6()),
		ConversationID: conversationId,
		AuthorID:       authorId,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Draft{
		Id:             draft.ID,
		ConversationId: draft.ConversationID,
		AuthorId:       draft.AuthorID,
		Text: domain.RichText{
			Text:  draft.Text,
			Parts: draft.TextParts,
		},
	}, nil
}

func (d *DraftRepository) Update(ctx context.Context, id uuid.UUID, updateFunc func(ctx context.Context, draft *domain.Draft) (*domain.Draft, error)) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := d.queries.WithTx(tx)
	rawDraft, err := queries.FindDraftForUpdate(ctx, id)
	if err != nil {
		return err
	}

	draft := d.mapToDomain(rawDraft)
	draft, err = updateFunc(ctx, draft)
	if err != nil {
		return err
	}

	err = queries.UpdateDraft(ctx, UpdateDraftParams{
		ID:        draft.Id,
		Text:      draft.Text.Text,
		TextParts: draft.Text.Parts,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *DraftRepository) RemoveByConversationAndAuthor(
	ctx context.Context,
	conversationId uuid.UUID,
	authorId uuid.UUID,
) error {
	return d.queries.RemoveDraft(ctx, RemoveDraftParams{
		ConversationID: conversationId,
		AuthorID:       authorId,
	})
}

func (d *DraftRepository) mapToDomain(draft MessengerDraft) *domain.Draft {
	return &domain.Draft{
		Id:             draft.ID,
		ConversationId: draft.ConversationID,
		AuthorId:       draft.AuthorID,
		Text: domain.RichText{
			Text:  draft.Text,
			Parts: draft.TextParts,
		},
	}
}
