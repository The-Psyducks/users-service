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

	// ModifyUser updates a user in the database
	ModifyUser(id uuid.UUID, data model.UpdateUserPrivateProfile) (model.UserRecord, error)

	// GetUserById retrieves a user from the database by its ID
	GetUserById(id uuid.UUID) (model.UserRecord, error)

	// GetUserByEmail retrieves a user from the database by its username
	// it is case sensitive
	GetUserByEmail(email string) (model.UserRecord, error)

	// CheckIfUsernameExists checks if a username already exists in the database
	// it is case insensitive
	CheckIfUsernameExists(username string) (bool, error)

	// CheckIfEmailExists checks if a mail already exists in the database
	// it is case insensitive
	CheckIfEmailExists(email string) (bool, error)

	// FollowUser associates a follower to a following user
	FollowUser(followerId uuid.UUID, followingId uuid.UUID) error

	// UnfollowUser removes a follower from a following user
	UnfollowUser(followerId uuid.UUID, followingId uuid.UUID) error

	// CheckIfUserFollows checks if followerID follows followingId
	CheckIfUserFollows(followerId uuid.UUID, followingId uuid.UUID) (bool, error)

	// GetAmountOfFollowers retrieves the amount of followers for a given user ID
	GetAmountOfFollowers(userId uuid.UUID) (int, error)

	// GetAmountOfFollowing retrieves the amount of following for a given user ID
	GetAmountOfFollowing(userId uuid.UUID) (int, error)

	// GetFollowers returns the followers for a given user ID and if there are more followers to retrieve
	// it also receives a timestamp, skip and limit to paginate the results
	GetFollowers(userId uuid.UUID, timestamp string, skip int, limit int) ([]model.UserRecord, bool, error)

	// GetAllUsers retrieves all the users in the database
	// it also receives a timestamp, skip and limit to paginate the results
	GetAllUsers(timestamp string, skip int, limit int) ([]model.UserRecord, bool, error)

	// GetFollowing returns the users that a user is following for a given user ID
	// and if there are more followers to retrieve.
	// It also receives a timestamp, skip and limit to paginate the results
	GetFollowing(userId uuid.UUID, timestamp string, skip int, limit int) ([]model.UserRecord, bool, error)

	// GetUsersWithUsernameContaining returns the users that have a username containing the text
	// it also receives a timestamp, skip and limit to paginate the results
	GetUsersWithUsernameContaining(text string, timestamp string, skip int, limit int) ([]model.UserRecord, bool, error)

	// GetAmountOfUsersWithUsernameContaining returns the amount of users that have a username containing the text
	GetAmountOfUsersWithUsernameContaining(text string) (int, error)

	// GetUsersWithOnlyNameContaining returns the users that JUST have the name containing the text. 
	// If the username also has it, it discards it
	// it also receives a timestamp, skip and limit to paginate the results
	GetUsersWithOnlyNameContaining(text string, timestamp string, skip int, limit int) ([]model.UserRecord, bool, error)

	// GetRecommendations returns the users that are recommended for a given user ID and if there are more users to retrieve
	// it also receives a timestamp, skip and limit to paginate the results
	// it calculates the recommendations based on the user's interests and location, returning first the users that share both
	// then the users that share only one of them
	GetRecommendations(userId uuid.UUID, timestamp string, skip int, limit int) ([]model.UserRecord, bool, error)

	// BlockUser blocks a user
	BlockUser(userId uuid.UUID, reason string) error

	// UnblockUser unblocks a user
	UnblockUser(userId uuid.UUID) error

	// CheckIfUserIsBlocked checks if a user is blocked
	CheckIfUserIsBlocked(userId uuid.UUID) (bool, error)

	// RegisterLoginAttempt registers a login attempt in the database
	RegisterLoginAttempt(userID uuid.UUID, provider *string, successful bool) error

	// GetLoginSummaryMetrics retrieves the login metrics
	GetLoginSummaryMetrics() (*model.LoginSummaryMetrics, error)

	// GetLocationMetrics retrieves the location metrics
	GetLocationMetrics() (*model.LocationMetrics, error)

	// GetUserBlockedMetrics retrieves the user blocked metrics
	GetUsersBlockedMetrics() (*model.UsersBlockedMetrics, error)
}
