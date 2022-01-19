package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mostafah/go-jalali/jalali"
	"gorm.io/gorm"
	"strconv"
	"time"
)

const (
	LAUNCH = "lunch"
	BUFFET = "buffet"
)

func (f *Food) AfterFind(tx *gorm.DB) (err error) {
	image := new(Image)
	image_count := tx.Where("id = ?", f.ImageID).Find(&image).RowsAffected
	if image_count > 0 {
		f.Image = *image
	}
	return
}

func (r *Reserve) BeforeUpdate(tx *gorm.DB) (err error) {
	prv_reserve := Reserve{}
	tx.Where("id = ?", r.ID).Find(&prv_reserve)

	food := Food{}
	tx.Where("id = ?", r.FoodID).Find(&food)
	y, m, d := jalali.Gtoj(food.Expire)
	jalali := fmt.Sprintf("%d-%d-%d", d, m, y)
	sign := r.Count - prv_reserve.Count
	if sign == 0 {
		return
	}
	payment := Payment{
		Amount:      (sign) * food.Price,
		UserID:      r.UserID,
		Description: food.Name + " (" + jalali + ") - (" + strconv.Itoa(sign) + "عدد) ",
	}
	err = tx.Table("payments").Create(&payment).Error

	if err != nil {
		return err
	}

	return
}

type (
	Food struct {
		Base
		Name    string     `json:"name"`
		Price   int        `json:"price"`
		Expire  time.Time  `json:"expire"`
		Type    string     `json:"type"`
		ImageID *uuid.UUID `json:"image_id"`
		Image   Image      `json:"image"`
	}

	Reserve struct {
		Base
		UserID uuid.UUID `json:"user_id"`
		FoodID uuid.UUID `json:"food_id"`
		Count  int       `json:"count"`
	}
)
