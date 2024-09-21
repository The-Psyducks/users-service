package model

import (
	"time"

	"github.com/google/uuid"
)

// UserProfileResponse is a struct that represents a user profile in the HTTP response
type UserProfileResponse struct {
	OwnProfile bool			`json:"own_profile" binding:"required"`
	Follows    bool			`json:"follows" binding:"required"`
	Profile    interface{}	`json:"profile" binding:"required"`
}

// UserPrivateProfile is a struct that represents a user in the HTTP response
type UserPrivateProfile struct {
	Id        uuid.UUID `json:"id" binding:"required"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	UserName  string    `json:"username" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Location  string    `json:"location" binding:"required"`
	Interests []string  `json:"interests" binding:"required"`
	Followers int       `json:"followers" binding:"required"`
	Following int       `json:"following" binding:"required"`
}

// UserPrivateProfile is a struct that represents a user in the HTTP response
type UserPublicProfile struct {
	Id        uuid.UUID `json:"id" binding:"required"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	UserName  string    `json:"username" binding:"required"`
	Location  string    `json:"location" binding:"required"`
	Followers int       `json:"followers" binding:"required"`
	Following int       `json:"following" binding:"required"`
}

// UserRecord is a struct that represents a user in the database
type UserRecord struct {
	Id        uuid.UUID `json:"id" db:"id"`
	UserName  string    `json:"username" db:"username"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Location  string    `json:"location" db:"location"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
