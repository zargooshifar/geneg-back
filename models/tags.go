package models

import "github.com/google/uuid"

type Tag struct {
	Base
	UserID uuid.UUID `json:"user_id"`
	TagID  string    `json:"tag_id"`
	Name   string    `json:"name"`
}
