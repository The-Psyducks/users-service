package service

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrRegistryNotFound = errors.New("registry not found")
)

const (
	IncorrectUsernameOrPassword = "Incorrect username or password"
	InternalServerError         = "Internal server error"
	EmailAlreadyExists 			= "Email already exists"
	UsernameAlreadyExists 		= "Username already exists"
	UsernameNotFound            = "User not found"
	RegistryNotFound            = "Registry not found"
	InvalidRegistryStep			= "Invalid registry step"
	InvalidInterest				= "Invalid interest"
	CantFollowYourself			= "Can't follow yourself"
	AlreadyFollowing			= "The user already follows this user"
	NotFollowing				= "The user is not following this user"
)
