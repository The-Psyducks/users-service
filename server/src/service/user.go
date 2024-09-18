package service

import (
	"users-service/src/database/registry_db"
	"users-service/src/database/users_db"

)

type User struct {
	userDb        users_db.UserDatabase
	registryDb    registry_db.RegistryDatabase
	userValidator *UserCreationValidator
}

func CreateUserService(userDb users_db.UserDatabase, registryDb registry_db.RegistryDatabase) *User {
	return &User{
		userDb:        userDb,
		registryDb:    registryDb,
		userValidator: NewUserCreationValidator(userDb),
	}
}
