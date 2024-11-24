package service

import (
	"users-service/src/database/registry_db"
	"users-service/src/database/users_db"
	amqp "github.com/rabbitmq/amqp091-go"
)

type User struct {
	userDb        users_db.UserDatabase
	registryDb    registry_db.RegistryDatabase
	userValidator *UserValidator
	amqpQueue	 *amqp.Channel	
}

func CreateUserService(userDb users_db.UserDatabase, registryDb registry_db.RegistryDatabase, queue *amqp.Channel) *User {
	return &User{
		userDb:        userDb,
		registryDb:    registryDb,
		userValidator: NewUserValidator(userDb),
		amqpQueue:     queue,
	}
}
