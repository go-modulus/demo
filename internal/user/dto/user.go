package dto

import (
	"time"
)

type User struct {
	Id           string `gorm:"primarykey"`
	Name         string
	Email        string
	RegisteredAt time.Time `gorm:"column:registered_at"`
}
