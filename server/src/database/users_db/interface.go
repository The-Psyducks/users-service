package users_db

import (
	"users-service/src/model"
	"github.com/google/uuid"
)

// UserDatabase interface to interact with the user's database
// it is used by the service layer
type UserDatabase interface {
	// CreateUser creates a new user in the database
	CreateUser(data model.UserRecord) (model.UserRecord, error)

	// GetUserById retrieves a user from the database by its ID
	GetUserById(id uuid.UUID) (model.UserRecord, error)

	// GetUserByUsername retrieves a user from the database by its username
	// it is case sensitive
	GetUserByUsername(username string) (model.UserRecord, error)

	// GetUserByEmail retrieves a user from the database by its username
	// it is case sensitive
	GetUserByEmail(email string) (model.UserRecord, error)

	// CheckIfUsernameExists checks if a username already exists in the database
	// it is case insensitive
	CheckIfUsernameExists(username string) (bool, error)

	// CheckIfEmailExists checks if a mail already exists in the database
	// it is case insensitive
	CheckIfEmailExists(email string) (bool, error)
}
