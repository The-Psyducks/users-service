package service

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
)

var (
	IncorrectUsernameOrPassword = "Incorrect username or password"
	InternalServerError         = "Internal server error"
	UsernameOrMailAlreadyExists = "Username or mail already exists"
	UsernameNotFound            = "User not found"
)
