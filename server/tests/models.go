package tests

import (
	"github.com/google/uuid"
)

type User struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Mail string `json:"mail" binding:"required"`
	Location int32 `json:"location" binding:"required"`
	Interests []int32 `json:"interests_ids" binding:"required"`
}

type UserResponse struct {
	Id   uuid.UUID    `json:"id" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Mail string `json:"mail" binding:"required"`
	Location string `json:"location" binding:"required"`
	Interests []string `json:"interests" binding:"required"`
}

type ErrorResponse struct {
    Type     string `json:"type"`
    Title    string `json:"title"`
    Status   int    `json:"status"`
    Detail   string `json:"detail"`
    Instance string `json:"instance"`
}