package database

import (
	"users-service/src/model"

	"github.com/google/uuid"
)

// UserDatabase interface to interact with the user's database
// it is used by the service layer
type UserDatabase interface {
	CreateUser(data model.UserRecord) (model.UserRecord, error)
	GetUserById(id string) (model.UserRecord, error)
}

// Database interface to interact with the Interest's database
// it is used by the service layer
type InterestsDatabase interface {
	AssociateInterestsToUser(userId uuid.UUID, interests []int32) ([]model.InterestRecord, error)
	GetInterestsNamesForUserId(id uuid.UUID) ([]string, error)
}
