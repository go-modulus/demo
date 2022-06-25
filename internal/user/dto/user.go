package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type User struct {
	Id           string `gorm:"primarykey"`
	Name         string
	Email        string
	RegisteredAt time.Time `gorm:"column:registered_at"`
}

func IdRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("Id is required."),
		is.UUID.Error("Id is not valid UUID4."),
	}
}

func EmailRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("Required field is empty."),
		is.EmailFormat.Error("Please enter a valid email address."),
	}
}

func NameRules() []validation.Rule {
	return []validation.Rule{
		validation.Required.Error("Required field is empty."),
		validation.Length(3, 50).Error("User name should be between 3 and 50 characters length."),
		is.Alpha.Error("User name should contain characters only."),
	}
}
