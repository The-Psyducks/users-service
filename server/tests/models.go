package tests

import (
	"github.com/google/uuid"
)

type ResolverResponse struct {
	NextAuthStep string      `json:"next_auth_step"`
	Metadata     interface{} `json:"metadata"`
}

type ResolverSignUpResponse struct {
	NextAuthStep string                 `json:"next_auth_step"`
	Metadata     MetadataSingUpResponse `json:"metadata"`
}

type MetadataSingUpResponse struct {
	OnboardingStep string `json:"onboarding_step"`
	RegistrationId string `json:"registration_id"`
}

type UserPersonalInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Location  int    `json:"location"`
}

type UserProfileResponse struct {
	OwnProfile bool        `json:"own_profile" binding:"required"`
	Follows    bool        `json:"follows" binding:"required"`
	Profile    interface{} `json:"profile" binding:"required"`
}

type UserPublicProfile struct {
	Id        uuid.UUID `json:"id" binding:"required"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	UserName  string    `json:"username" binding:"required"`
	Location  string    `json:"location" binding:"required"`
	Followers int       `json:"followers" binding:"required"`
	Following int       `json:"following" binding:"required"`
}

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

type Location struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Interest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RegisterOptions struct {
	Locations []Location `json:"locations"`
	Interests []Interest `json:"interests"`
}

type ErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Type     string            `json:"type"`
	Title    string            `json:"title"`
	Status   int               `json:"status"`
	Detail   string            `json:"detail"`
	Instance string            `json:"instance"`
	Errors   []ValidationError `json:"errors"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Profile	 UserPrivateProfile `json:"profile"`
}

type EditUserProfileRequest struct {
	Username string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Location  int    `json:"location"`
	Interests []int  `json:"interests"`
}

type FollowUserProfile struct {
	Follows bool `json:"follows"`
	Profile UserPublicProfile `json:"profile"`
}

type Pagination struct {
	NextOffset int `json:"next_offset"`
	Limit      int `json:"limit"`
}

type FollowersResponse struct {
	Followers []FollowUserProfile	`json:"data"`
	Pagination Pagination 			`json:"pagination"`
}

type FollowingResponse struct {
	Following []FollowUserProfile	`json:"data"`
	Pagination Pagination			`json:"pagination"`
}

