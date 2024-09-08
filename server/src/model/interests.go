package model

import "github.com/google/uuid"

type Interest struct {
	Id       int    `json:"id"`
	Interest string `json:"name"`
}

type InterestRecord struct {
	InterestId int       `json:"id" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	UserId     uuid.UUID `json:"user_id" binding:"required"`
}
