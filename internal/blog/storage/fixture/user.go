package fixture

import (
	"boilerplate/internal/blog/storage"
	"context"
	"github.com/gofrs/uuid"
)

type PostFixture struct {
	blogDb *storage.Queries
}

func NewPostFixture(blogDb *storage.Queries) *PostFixture {
	return &PostFixture{
		blogDb: blogDb,
	}
}

func (f *PostFixture) CreateRandomPost() (storage.Post, func(), string) {
	name := "test"
	id, _ := uuid.NewV6()

	return f.CreateParticularPost(id, name)
}

func (f *PostFixture) CreateParticularPost(
	authorId uuid.UUID,
	title string,

) (storage.Post, func(), string) {
	id, _ := uuid.NewV6()
	post, _ := f.blogDb.CreatePost(
		context.Background(), storage.CreatePostParams{
			ID:       id,
			Title:    "",
			Body:     "",
			AuthorID: authorId,
			Slug:     "",
		},
	)
	return post, func() {
		f.DeletePost(id)
	}, "The post " + title

}

func (f *PostFixture) DeletePost(id uuid.UUID) {
	_ = f.blogDb.DeletePost(context.Background(), id)
}
