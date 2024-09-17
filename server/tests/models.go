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

type UserProfile struct {
	Id        uuid.UUID `json:"id" binding:"required"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	UserName  string    `json:"username" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Location  string    `json:"location" binding:"required"`
	Interests []string  `json:"interests" binding:"required"`
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
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
