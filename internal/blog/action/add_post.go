package action

import (
	"boilerplate/internal/blog/storage"
	"boilerplate/internal/framework"
	validator "boilerplate/internal/ozzo-validator"
	"context"
	"github.com/ggicci/httpin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
	"regexp"
)

var slugExcludeRegexp = regexp.MustCompile(`[^a-z0-9\-_]+`)

type AddPostRequest struct {
	httpin.BodyDecoder
	Title              string `json:"title"`
	Body               string `json:"body"`
	PublishImmediately bool   `json:"publish_immediately"`
}

type AddPostAction struct {
	blogDb *storage.Queries
}

func NewAddPostAction(blogDb *storage.Queries) *AddPostAction {
	return &AddPostAction{blogDb: blogDb}
}

func InitAddPostAction(
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
	action *AddPostAction,
) error {
	handler, err := framework.WrapHandler[*AddPostRequest, storage.Post](errorHandler, action, 201)

	if err != nil {
		return err
	}
	routes.Post("/api/blog/posts", handler)

	return nil
}

func (r *AddPostRequest) Validate(ctx context.Context) *framework.ValidationErrors {
	err := validation.ValidateStructWithContext(
		ctx,
		r,
		validation.Field(
			&r.Title,
			validation.Required.Error("Title is required"),
			is.Alpha.Error("Title should be alphabetical."),
			validation.Length(3, 100).Error("Title should be from 3 to 100 characters length."),
		),
		validation.Field(
			&r.Body,
			validation.Required.Error("Body is required"),
			validation.Length(1, 1000).Error("Post text should be from 1 to 1000 characters length."),
		),
	)

	return validator.AsAppValidationErrors(err)
}

func (a *AddPostAction) Handle(ctx context.Context, req *AddPostRequest) (storage.Post, error) {
	id := uuid.Must(uuid.NewV6())
	slug := slugExcludeRegexp.ReplaceAllString(req.Title, "_")

	params := storage.CreatePostParams{
		ID:       id,
		Title:    req.Title,
		Body:     req.Body,
		AuthorID: uuid.Must(uuid.NewV6()),
		Slug:     slug,
	}
	result, err := a.blogDb.CreatePost(ctx, params)
	if err != nil {
		return storage.Post{}, err
	}

	if req.PublishImmediately {
		published, err := a.blogDb.PublishPost(ctx, id)
		if err != nil {
			return result, err
		}
		result = published
	}

	return result, nil
}
