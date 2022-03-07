package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (u *CheckIn) AfterFind(tx *gorm.DB) (err error) {
	if u.UserID != nil {
		user := new(User)
		tx.Where("id = ?", u.UserID).Find(&user)
		u.User = *user
	}
	return
}

type (
	CheckIn struct {
		Base
		UserID *uuid.UUID `json:"user_id"`
		User   User       `json:"user"`
		Tagged bool       `json:"tagged"`
	}
)
