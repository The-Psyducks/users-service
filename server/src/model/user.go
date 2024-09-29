package model

import (
	"time"

	"github.com/google/uuid"
)

// UserProfileResponse is a struct that represents a user profile in the HTTP response
type UserProfileResponse struct {
	OwnProfile bool        `json:"own_profile" binding:"required"`
	Follows    bool        `json:"follows" binding:"required"`
	Profile    interface{} `json:"profile" binding:"required"`
}

// UserPrivateProfileRequest is a struct that represents a user in the HTTP request
type UpdateUserPrivateProfile struct {
	UserName    string `json:"username" binding:"required"`
	PicturePath string `json:"picture_path" db:"picture_path"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Interests   []string  `json:"interests" binding:"required"`
}

type UpdateUserPrivateProfileRequest struct {
	PicturePath string `json:"picture_path"`
    FirstName  string `json:"first_name" validate:"firstnamevalidator"`
    LastName   string `json:"last_name" validate:"lastnamevalidator"`
    UserName   string `json:"username" validate:"usernamevalidator"`
	Location    int    `json:"location" validate:"locationvalidator"`
	Interests   []int  `json:"interests" validate:"interestsvalidator"`
}

// UserPrivateProfile is a struct that represents a user in the HTTP response
type UserPrivateProfile struct {
	Id          uuid.UUID `json:"id" binding:"required"`
	UserName    string    `json:"username" binding:"required"`
	PicturePath string    `json:"picture_path" db:"picture_path"`
	FirstName   string    `json:"first_name" binding:"required"`
	LastName    string    `json:"last_name" binding:"required"`
	Email       string    `json:"email" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Interests   []string  `json:"interests" binding:"required"`
	Followers   int       `json:"followers" binding:"required"`
	Following   int       `json:"following" binding:"required"`
}

// UserPrivateProfile is a struct that represents a user in the HTTP response
type UserPublicProfile struct {
	Id          uuid.UUID `json:"id" binding:"required"`
	UserName    string    `json:"username" binding:"required"`
	PicturePath string    `json:"picture_path" db:"picture_path"`
	FirstName   string    `json:"first_name" binding:"required"`
	LastName    string    `json:"last_name" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Followers   int       `json:"followers" binding:"required"`
	Following   int       `json:"following" binding:"required"`
}

// UserRecord is a struct that represents a user in the database
type UserRecord struct {
	Id          uuid.UUID `json:"id" db:"id"`
	UserName    string    `json:"username" db:"username"`
	PicturePath string    `json:"picture_path" db:"picture_path"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"password" db:"password"`
	Location    string    `json:"location" db:"location"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
