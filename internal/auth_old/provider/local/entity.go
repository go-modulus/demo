package local

import (
	"time"
)

type LocalAccount struct {
	UserID    string  `gorm:"column:user_id"`
	Email     *string `gorm:"column:email"`
	Nickname  *string `gorm:"column:nickname"`
	Phone     *string `gorm:"column:phone"`
	Password  string  `gorm:"column:password"`
	CreatedAt time.Time
}

func (a LocalAccount) GetUserId() string {
	return a.UserID
}
