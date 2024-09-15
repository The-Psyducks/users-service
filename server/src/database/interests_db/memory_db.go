package interests_db

import (
	"users-service/src/database"

	"github.com/google/uuid"
)

// InterestsMemoryDB struct to hold the data
type InterestsMemoryDB struct {
	userInterests map[uuid.UUID]map[string]bool
}

// CreateInterestsMemoryDB creates a new instance of InMemoryDB
func CreateInterestsMemoryDB() (*InterestsMemoryDB, error) {
	return &InterestsMemoryDB{
		userInterests: make(map[uuid.UUID]map[string]bool),
	}, nil
}

// AssociateInterestsToUser associates interests to a user
func (db *InterestsMemoryDB) AssociateInterestsToUser(userId uuid.UUID, interests []string) error {
	if _, exists := db.userInterests[userId]; !exists {
		db.userInterests[userId] = make(map[string]bool)
	}

	for _, interest := range interests {
		if _, exists := db.userInterests[userId][interest]; exists {
			return database.ErrKeyAlreadyExists
		}
		db.userInterests[userId][interest] = true
	}
	return nil
}

// GetInterestsForUserId retrieves interests for a given user ID
func (db *InterestsMemoryDB) GetInterestsForUserId(id uuid.UUID) ([]string, error) {
	interestSet, exists := db.userInterests[id]
	if !exists {
		return nil, database.ErrKeyNotFound
	}

	interests := make([]string, 0, len(interestSet))
	for interest := range interestSet {
		interests = append(interests, interest)
	}
	return interests, nil
}
