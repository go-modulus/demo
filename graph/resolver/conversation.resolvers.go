package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.22

import (
	"context"
	"demo/graph/generated"
	"demo/graph/model"

	"github.com/gofrs/uuid"
)

// CreateOneToOneConversation is the resolver for the createOneToOneConversation field.
func (r *mutationResolver) CreateOneToOneConversation(ctx context.Context, receiverID uuid.UUID) (model.CreateOneToOneConversationResult, error) {
	return r.messengerResolver.CreateOneToOneConversation(ctx, receiverID)
}

// Draft is the resolver for the draft field.
func (r *oneToOneConversationResolver) Draft(ctx context.Context, obj *model.OneToOneConversation) (*model.Draft, error) {
	return r.messengerResolver.Draft(ctx, obj)
}

// LastMessage is the resolver for the lastMessage field.
func (r *oneToOneConversationResolver) LastMessage(ctx context.Context, obj *model.OneToOneConversation) (model.Message, error) {
	return r.messengerResolver.LastMessage(ctx, obj)
}

// Conversation is the resolver for the conversation field.
func (r *queryResolver) Conversation(ctx context.Context, id uuid.UUID) (model.ConversationResult, error) {
	return r.messengerResolver.GetConversation(ctx, id)
}

// MyConversations is the resolver for the myConversations field.
func (r *queryResolver) MyConversations(ctx context.Context, first int, after *string) (model.MyConversationsResult, error) {
	return r.messengerResolver.MyConversations(ctx, first, after)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// OneToOneConversation returns generated.OneToOneConversationResolver implementation.
func (r *Resolver) OneToOneConversation() generated.OneToOneConversationResolver {
	return &oneToOneConversationResolver{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type oneToOneConversationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }