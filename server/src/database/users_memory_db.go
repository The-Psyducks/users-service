package database

import (
	"sort"
	"time"
	"strings"
	"users-service/src/model"

	"github.com/google/uuid"
)

// UserMemoryDB is a simple in-memory database
type UserMemoryDB struct {
	data map[string]model.UserRecord
}

// NewUserMemoryDB creates a new MemoryDB
func NewUserMemoryDB() (*UserMemoryDB, error) {
	return &UserMemoryDB{
		data: make(map[string]model.UserRecord),
	}, nil
}

func (m *UserMemoryDB) CreateUser(data model.UserRecord) (model.UserRecord, error) {
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

	m.data[newUser.Id.String()] = newUser
	return newUser, nil
}

// GetUserById retrieves a user from the database by its ID
func (m *UserMemoryDB) GetUserById(id string) (model.UserRecord, error) {
	user, found := m.data[id]
	if !found {
		return model.UserRecord{}, ErrKeyNotFound
	}
	return user, nil
}

func (m *UserMemoryDB) CheckIfUsernameExists(username string) (bool, error) {
	for _, user := range m.data {
		if strings.EqualFold(user.UserName, username) {
			return true, nil
		}
	}
	return false, nil
}

func (m *UserMemoryDB) CheckIfMailExists(mail string) (bool, error) {
	for _, user := range m.data {
		if strings.EqualFold(user.Mail, mail) {
			return true, nil
		}
	}
	return false, nil
}

func (m *UserMemoryDB) GetUserByUsername(username string) (model.UserRecord, error) {
	for _, user := range m.data {
		if user.UserName == username {
			return user, nil
		}
	}
	return model.UserRecord{}, ErrKeyNotFound
}

func (m *UserMemoryDB) Delete(id string) error {
	_, found := m.data[id]
	if !found {
		return ErrKeyNotFound
	}
	delete(m.data, id)
	return nil
}

func (m *UserMemoryDB) GetAll() ([]model.UserRecord, error) {
	allUsers := make([]model.UserRecord, 0, len(m.data))
	for _, user := range m.data {
		allUsers = append(allUsers, user)
	}
	return allUsers, nil
}

func (m *UserMemoryDB) GetAllInOrder() ([]model.UserRecord, error) {
	allUsers, err := m.GetAll()
	if err != nil {
		return nil, err
	}
	sort.Slice(allUsers, func(i, j int) bool {
		return allUsers[i].CreatedAt.After(allUsers[j].CreatedAt)
	})
	return allUsers, nil
}
