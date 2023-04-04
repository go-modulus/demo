package fixture

import (
	"boilerplate/internal/user/storage"
	"context"
	"github.com/gofrs/uuid"
)

type UserFixture struct {
	userDb *storage.Queries
}

func NewUserFixture(userDb *storage.Queries) *UserFixture {
	return &UserFixture{
		userDb: userDb,
	}
}

func (f *UserFixture) CreateRandomUser() (storage.User, func(), string) {
	name := "test"
	id, _ := uuid.NewV6()

	return f.CreateParticularUser(id, name)
}

func (f *UserFixture) CreateParticularUser(id uuid.UUID, name string) (storage.User, func(), string) {
	user, _ := f.userDb.CreateUser(
		context.Background(), storage.CreateUserParams{
			ID:    id,
			Name:  name,
			Email: name + "+test" + id.String() + "@gmail.com",
		},
	)
	return user, func() {
		f.DeleteUser(id)
	}, "The user " + id.String()

}

func (f *UserFixture) DeleteUser(id uuid.UUID) {
	_ = f.userDb.DeleteUser(context.Background(), id)
}
