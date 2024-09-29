package users_db

import (
	"sort"
	"strings"
	"time"
	"users-service/src/database"
	"users-service/src/model"

	"github.com/google/uuid"
)

// UserMemoryDB is a simple in-memory database
type UserMemoryDB struct {
	data map[string]model.UserRecord
	userInterests map[uuid.UUID]map[string]bool
}

// CreateUserMemoryDB creates a new MemoryDB
func CreateUserMemoryDB() (*UserMemoryDB, error) {
	return &UserMemoryDB{
		data: make(map[string]model.UserRecord),
		userInterests: make(map[uuid.UUID]map[string]bool),
	}, nil
}

func (m *UserMemoryDB) CreateUser(data model.UserRecord) (model.UserRecord, error) {
	newUser := model.UserRecord{
		Id:        uuid.New(),
		UserName:  data.UserName,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
		Location:  data.Location,
		CreatedAt: time.Now().UTC(),
	}

	m.data[newUser.Id.String()] = newUser
	return newUser, nil
}

// GetUserById retrieves a user from the database by its ID
func (m *UserMemoryDB) GetUserById(id uuid.UUID) (model.UserRecord, error) {
	user, found := m.data[id.String()]
	if !found {
		return model.UserRecord{}, database.ErrKeyNotFound
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

func (m *UserMemoryDB) CheckIfEmailExists(email string) (bool, error) {
	for _, user := range m.data {
		if strings.EqualFold(user.Email, email) {
			return true, nil
		}
	}
	return false, nil
}

func (m *UserMemoryDB) GetUserByEmail(email string) (model.UserRecord, error) {
	for _, user := range m.data {
		if user.Email == email {
			return user, nil
		}
	}
	return model.UserRecord{}, database.ErrKeyNotFound
}

func (m *UserMemoryDB) Delete(id uuid.UUID) error {
	_, found := m.data[id.String()]
	if !found {
		return database.ErrKeyNotFound
	}
	delete(m.data, id.String())
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

// AssociateInterestsToUser associates interests to a user
func (db *UserMemoryDB) AssociateInterestsToUser(userId uuid.UUID, interests []string) error {
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
func (db *UserMemoryDB) GetInterestsForUserId(id uuid.UUID) ([]string, error) {
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