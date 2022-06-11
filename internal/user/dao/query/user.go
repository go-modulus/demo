package query

import (
	"context"
	"gorm.io/gorm"
)

const UserTable = `"user"."user"`

type UserQuery struct {
	Db *gorm.DB
}

func NewUserQuery(ctx context.Context, db *gorm.DB) *UserQuery {
	localCopy := db.Table(UserTable).WithContext(ctx)
	query := &UserQuery{
		Db: localCopy,
	}
	return query
}

func (p *UserQuery) Email(email string) *UserQuery {
	p.Db = p.Db.Where(UserTable+".email = ?", email)

	return p
}

func (p *UserQuery) Id(id string) *UserQuery {
	p.Db = p.Db.Where(UserTable+".id = ?", id)
	return p
}

func (p *UserQuery) NewerFirst() *UserQuery {
	p.Db = p.Db.Order(UserTable + ".registered_at DESC")
	return p
}

func (p *UserQuery) Build() *gorm.DB {
	//resultDb := p.Db.Session(&gorm.Session{})

	return p.Db
}
