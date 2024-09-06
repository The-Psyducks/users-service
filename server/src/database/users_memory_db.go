package database

import (
	"time"
	"users-service/src/model"

	"github.com/google/uuid"
)

// MemoryDB is a simple in-memory database
type MemoryDB struct {
	db genericDB[model.UserRecord]
}

// CreateUser creates a new user in the database
func (m *MemoryDB) CreateUser(data model.UserRecord) (model.UserRecord, error) {
	newUser := model.UserRecord{
		Id:        uuid.New(),
		UserName:  data.UserName,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Mail:      data.Mail,
		Password:  data.Password,
		Location:  data.Location,
		CreatedAt: time.Now().UTC(),
	}

	return m.db.Create(newUser, newUser.Id.String())
}

// GetUserById retrieves a user from the database by its ID
func (m *MemoryDB) GetUserById(id string) (model.UserRecord, error) {
	return m.db.Get(id)
}

// NewUserMemoryDB creates a new MemoryDB
func NewUserMemoryDB() (*MemoryDB, error) {
	return &MemoryDB{
		db: CreateDB[model.UserRecord](),
	}, nil
}
