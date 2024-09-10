package interests_db

import (
	"users-service/src/model"

	"github.com/google/uuid"
)
// Database interface to interact with the Interest's database
// it is used by the service layer
type InterestsDatabase interface {
	// AssociateInterestsToUser associates interests to a user
	AssociateInterestsToUser(userId uuid.UUID, interests []int) ([]model.InterestRecord, error)

	// GetInterestsForUserId retrieves interests for a given user ID
	GetInterestsNamesForUserId(id uuid.UUID) ([]string, error)
}
