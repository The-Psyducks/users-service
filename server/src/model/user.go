package model

import (
	"time"

	"github.com/google/uuid"
)

// UserRequest is a struct that represents a user HTTP request
type UserRequest struct {
	UserName string `json:"username" binding:"required"`
	Name string `json:"name" binding:"required"`
	Mail string `json:"mail" binding:"required"`
	Location string `json:"location" binding:"required"`
}

// UserResponse is a struct that represents a user in the HTTP response
type UserResponse struct {
	Id   uuid.UUID    `json:"id" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Name string `json:"name" binding:"required"`
	Mail string `json:"mail" binding:"required"`
	Location string `json:"location" binding:"required"`
}

// UserRecord is a struct that represents a user in the database
type UserRecord struct {
	Id   uuid.UUID    `json:"id" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Name string `json:"name" binding:"required"`
	Mail string `json:"mail" binding:"required"`
	Location string `json:"location" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (u UserRecord) GetCreatedAt() time.Time {
	return u.CreatedAt
}