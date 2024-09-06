package model

import (
	"time"

	"github.com/google/uuid"
)

// UserRequest is a struct that represents a user HTTP request
type UserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Mail string `json:"mail" binding:"required"`
	Location string `json:"location" binding:"required"`
	Interests []string `json:"interests" binding:"required"`
}

// UserResponse is a struct that represents a user in the HTTP response
type UserResponse struct {
	Id   uuid.UUID    `json:"id" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Mail string `json:"mail" binding:"required"`
	Location string `json:"location" binding:"required"`
	Interests []string `json:"interests" binding:"required"`
}

// UserRecord is a struct that represents a user in the database
type UserRecord struct {
	Id   uuid.UUID    `json:"id" binding:"required"`
	UserName string `json:"username" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Mail string `json:"mail" binding:"required"`
	Password string `json:"password" binding:"required"`
	Location string `json:"location" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (u UserRecord) GetCreatedAt() time.Time {
	return u.CreatedAt
}