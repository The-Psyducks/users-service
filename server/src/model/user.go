package model

import (
	"time"

	"github.com/google/uuid"
)

// UserRequest is a struct that represents a user HTTP request
type UserRequest struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	UserName     string `json:"username" validate:"required,usernamevalidator"`
	Password     string `json:"password" validate:"required,passwordvalidator"`
	Mail         string `json:"mail" validate:"required,mailvalidator"`
	LocationId   int    `json:"location" validate:"locationvalidator"`
	InterestsIds []int  `json:"interests_ids" validate:"required,interestsvalidator"`
}

// UserResponse is a struct that represents a user in the HTTP response
type UserResponse struct {
	Id        uuid.UUID `json:"id" binding:"required"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	UserName  string    `json:"username" binding:"required"`
	Mail      string    `json:"mail" binding:"required"`
	Location  string    `json:"location" binding:"required"`
	Interests []string  `json:"interests" binding:"required"`
}

// UserRecord is a struct that represents a user in the database
type UserRecord struct {
	Id        uuid.UUID `json:"id"`
	UserName  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Mail      string    `json:"mail"`
	Password  string    `json:"password"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}
