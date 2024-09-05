package database

import "users-service/src/model"

// Database interface to interact with the database
// it is used by the service layer to interact with the database
type Database interface {
	CreateUser(data model.UserRequest) (model.UserRecord, error)
	GetUserById(id string) (model.UserRecord, error)
}