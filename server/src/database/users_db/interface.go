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

	// AssociateInterestsToUser associates interests to a user
	AssociateInterestsToUser(userId uuid.UUID, interests []string) error

	// GetInterestsForUserId retrieves interests for a given user ID
	GetInterestsForUserId(id uuid.UUID) ([]string, error)

	// FollowUser associates a follower to a following user
	FollowUser(followerId uuid.UUID, followingId uuid.UUID) error

	// UnfollowUser removes a follower from a following user
	UnfollowUser(followerId uuid.UUID, followingId uuid.UUID) error

	// CheckIfUserFollows checks if a user follows another user
	CheckIfUserFollows(followerId string, followingId string) (bool, error)

	// GetAmountOfFollowers retrieves the amount of followers for a given user ID
	GetAmountOfFollowers(userId uuid.UUID) (int, error)

	// GetAmountOfFollowing retrieves the amount of following for a given user ID
	GetAmountOfFollowing(userId uuid.UUID) (int, error)

	// GetFollowers retrieves the followers for a given user ID
	GetFollowers(userId uuid.UUID) ([]model.UserRecord, error)

	// GetFollowing retrieves the following for a given user ID
	GetFollowing(userId uuid.UUID) ([]model.UserRecord, error)
}

