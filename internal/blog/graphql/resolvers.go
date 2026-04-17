package graphql

import (
	"context"
	"fmt"

	"github.com/go-modulus/auth"
	"github.com/go-modulus/demo/internal/blog/storage"
	"github.com/go-modulus/demo/internal/graphql/model"
	"github.com/go-modulus/modulus/errors"
	"github.com/go-modulus/modulus/validator"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofrs/uuid"
)

type Resolver struct {
	blogQueries *storage.Queries
}

func NewResolver(blogQueries *storage.Queries) *Resolver {
	return &Resolver{blogQueries: blogQueries}
}

func (r *Resolver) CreatePost(ctx context.Context, input model.CreatePostInput) (storage.Post, error) {
	// validate input using Ozzo validation wrapped by modulus validator
	err := validator.ValidateStructWithContext(
		ctx,
		&input,
		validation.Field(
			&input.Title,
			validation.Required.Error("Title is required"),
		),
		validation.Field(
			&input.Content,
			validation.Required.Error("Content is required"),
		),
	)
	if err != nil {
		return storage.Post{}, errors.WithTrace(err)
	}

	preview := input.Content
	if len(input.Content) > 100 {
		preview = input.Content[0:100]
	}

	authorID := auth.GetPerformerID(ctx)
	if authorID == uuid.Nil {
		return storage.Post{}, errors.WithTrace(auth.ErrUnauthenticated)
	}

	return r.blogQueries.CreatePost(
		ctx, storage.CreatePostParams{
			ID:       uuid.Must(uuid.NewV6()),
			Title:    input.Title,
			Preview:  preview,
			Content:  input.Content,
			AuthorID: authorID,
		},
	)
}

// PublishPost is the resolver for the publishPost field.
func (r *Resolver) PublishPost(ctx context.Context, id uuid.UUID) (storage.Post, error) {
	return r.blogQueries.PublishPost(ctx, id)
}

// DeletePost is the resolver for the deletePost field.
func (r *Resolver) DeletePost(ctx context.Context, id uuid.UUID) (bool, error) {
	panic(fmt.Errorf("not implemented: DeletePost - deletePost"))
}

// Post is the resolver for the post field.
func (r *Resolver) Post(ctx context.Context, id string) (*storage.Post, error) {
	panic(fmt.Errorf("not implemented: Post - post"))
}

// Posts is the resolver for the posts field.
func (r *Resolver) Posts(ctx context.Context) ([]storage.Post, error) {
	authorID := auth.GetPerformerID(ctx)
	if authorID == uuid.Nil {
		return nil, auth.ErrUnauthenticated
	}
	return r.blogQueries.FindPosts(ctx, authorID)
}
