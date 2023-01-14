package messenger

import (
	"database/sql"
	"demo/internal/messenger/api"
	"demo/internal/messenger/domain"
	"demo/internal/messenger/persistence"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"messenger",
		fx.Provide(
			api.NewResolver,
			api.NewTransformer,
			api.NewDraftLoaderFactory,
			api.NewLastMessageLoaderFactory,

			domain.NewConversationViewer,
			domain.NewConversationCreator,
			domain.NewConversationPolicy,

			domain.NewRichTextParser,

			domain.NewDraftSaver,
			domain.NewDraftRemover,

			domain.NewMessageViewer,
			domain.NewMessageCreator,
			domain.NewMessageEditor,
			domain.NewMessageRemover,
			domain.NewMessagePolicy,

			persistence.NewConversationRepository,
			persistence.NewDraftRepository,
			persistence.NewMessageRepository,

			func(repo *persistence.ConversationRepository) domain.ConversationRepository {
				return repo
			},
			func(repo *persistence.DraftRepository) domain.DraftRepository {
				return repo
			},
			func(repo *persistence.MessageRepository) domain.MessageRepository {
				return repo
			},
			func(db *sql.DB) persistence.DBTX {
				return db
			},
			func(db persistence.DBTX) *persistence.Queries {
				return persistence.New(db)
			},
		),
	)
}
