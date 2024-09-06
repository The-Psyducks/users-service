package model

import "github.com/google/uuid"

type InterestRequest struct {
	Id  string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	UserId uuid.UUID `json:"user_id" binding:"required"`
}