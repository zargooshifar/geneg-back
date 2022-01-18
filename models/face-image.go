package models

import "github.com/google/uuid"

type FaceImage struct {
	Base
	UserID uuid.UUID `json:"user_id"`
	DataURL string `json:"data_url"`
}

type FaceID struct {
	ID int `gorm:"primary_key;" json:"id"`
	UserID uuid.UUID
}
