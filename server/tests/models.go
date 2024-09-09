package tests

import (
	"github.com/google/uuid"
)

type User struct {
    FirstName string  	`json:"first_name"`
    LastName  string  	`json:"last_name"`
    UserName  string  	`json:"username"`
    Password  string  	`json:"password"`
    Mail      string  	`json:"mail"`
    Location  int	  	`json:"location"`
    Interests []int		`json:"interests_ids"`
}

type UserProfile struct {
    Id        uuid.UUID `json:"id" binding:"required"`
    FirstName string    `json:"first_name" binding:"required"`
    LastName  string    `json:"last_name" binding:"required"`
    UserName  string    `json:"username" binding:"required"`
    Mail      string    `json:"mail" binding:"required"`
    Location  string    `json:"location" binding:"required"`
    Interests []string  `json:"interests" binding:"required"`
}

type Location struct {
	Id int `json:"id"`
	Name string `json:"name"`
} 

type Interest struct {
	Id int		`json:"id"`
	Name string	`json:"name"`
}

type RegisterOptions struct {
	Locations []Location	`json:"locations"`
	Interests []Interest 	`json:"interests"`
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
	Valid bool `json:"valid"`
}