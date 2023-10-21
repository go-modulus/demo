package action

import (
	"boilerplate/internal/auth"
	"boilerplate/internal/blog/storage"
	"boilerplate/internal/framework"
	validator "boilerplate/internal/ozzo-validator"
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetPostsRequest struct {
	Count int `in:"query=posts_count;default=10"`
	Page  int `in:"query=posts_page;default=1"`
}

func (r *GetPostsRequest) Validate(ctx context.Context) *framework.ValidationErrors {
	err := validation.ValidateStructWithContext(
		ctx,
		r,
		validation.Field(
			&r.Count,
			validation.Required.Error("Count is required"),
			validation.Min(1).Error("Count should be more than 0."),
			validation.Max(10).Error("Count should be less than or equal 10."),
		),
	)

	return validator.AsAppValidationErrors(err)
}

type Pagination struct {
	PagesCount int `json:"pages_count"`
	Page       int `json:"page"`
}

type PostsResponse struct {
	Data       []storage.Post `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

type GetPostsAction struct {
	blogDb *storage.Queries
}

func NewGetPostsAction(blogDb *storage.Queries) *GetPostsAction {
	return &GetPostsAction{blogDb: blogDb}
}

func InitGetPostsAction(
	auth *auth.Auth,
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
	action *GetPostsAction,
) error {
	getPosts, err := framework.WrapHandler[*GetPostsRequest, PostsResponse](errorHandler, action, 200)

	if err != nil {
		return err
	}
	routes.Get("/api/blog/posts", getPosts)

	return nil
}

func (a *GetPostsAction) Handle(ctx context.Context, req *GetPostsRequest) (PostsResponse, error) {
	if err := req.Validate(ctx); err != nil {
		return PostsResponse{}, err
	}
	data, err := a.blogDb.ListPosts(
		ctx, storage.ListPostsParams{
			After: int32(req.Page-1) * int32(req.Count),
			Count: int32(req.Count),
		},
	)

	if err != nil {
		return PostsResponse{}, err
	}
	count, err := a.blogDb.CountPosts(ctx)
	if err != nil {
		return PostsResponse{}, err
	}
	pagesCount := (int(count) / int(req.Count)) + 1
	r := PostsResponse{
		Data:       data,
		Pagination: Pagination{PagesCount: pagesCount, Page: req.Page},
	}
	return r, nil
}
