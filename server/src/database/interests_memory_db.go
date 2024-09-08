package database

import (
	"users-service/src/model"
	"github.com/google/uuid"
)

// InMemoryDB struct to hold the data
type InMemoryDB struct {
    userInterests map[uuid.UUID]map[int]bool // userId -> map of interestId -> bool
}

// NewInterestsMemoryDB creates a new instance of InMemoryDB
func NewInterestsMemoryDB() (*InMemoryDB, error) {
    return &InMemoryDB{
        userInterests: make(map[uuid.UUID]map[int]bool),
    }, nil
}

// AssociateInterestsToUser associates interests to a user
func (db *InMemoryDB) AssociateInterestsToUser(userId uuid.UUID, interests []int) ([]model.InterestRecord, error) {
    if _, exists := db.userInterests[userId]; !exists {
        db.userInterests[userId] = make(map[int]bool)
    }

    var records []model.InterestRecord

    for _, interestId := range interests {
        if interestName, exists := predefinedInterests[interestId]; exists {
            db.userInterests[userId][interestId] = true
            records = append(records, model.InterestRecord{
                InterestId: interestId,
                Name:       interestName,
                UserId:     userId,
            })
        }
    }

    return records, nil
}

// GetInterestsForUserId retrieves interests for a given user ID
func (db *InMemoryDB) GetInterestsNamesForUserId(id uuid.UUID) ([]string, error) {
    var records []string

    if userInterests, exists := db.userInterests[id]; exists {
        for interestId := range userInterests {
            if interestName, exists := predefinedInterests[interestId]; exists {
                records = append(records, interestName)
            }
        }
    }

    return records, nil
}