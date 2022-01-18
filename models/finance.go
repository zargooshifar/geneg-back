package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (p *Payment) AfterSave(tx *gorm.DB) (err error) {
	user := User{}
	tx.Where("id = ?", p.UserID).Find(&user)
	user.Balance += p.Amount
	tx.Save(&user)

	return
}

type (
	Payment struct {
		Base
		gorm.Model
		Amount int `json:"amount"`
		Description string `json:"description"`
		UserID	uuid.UUID `json:"user_id"`
	}
)
