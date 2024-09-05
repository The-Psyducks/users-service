package database

import (
	"users-service/src/model"
	"time"
	"github.com/google/uuid"
)

// MemoryDB is a simple in-memory database
type MemoryDB struct {
	db genericDB[model.UserRecord]
}

// CreateUser creates a new user in the database
func (m *MemoryDB) CreateUser(data model.UserRequest) (model.UserRecord, error) {
	newUser := model.UserRecord{
		Id:   uuid.New(),
		UserName: data.UserName,
		Name: data.Name,
		Mail: data.Mail,
		Location: data.Location,
		CreatedAt: time.Now().UTC(),
	}

	return m.db.Create(newUser, newUser.Id.String())
}

// GetUserById retrieves a user from the database by its ID
func (m *MemoryDB) GetUserById(id string) (model.UserRecord, error) {
	return m.db.Get(id)
}

// NewMemoryDB creates a new MemoryDB
func NewMemoryDB() (*MemoryDB, error) {
	return &MemoryDB{
		db: CreateDB[model.UserRecord](),
	}, nil
}