package registry_db

import (
	"log/slog"
	"users-service/src/database"
	"users-service/src/model"

	"github.com/google/uuid"
)

type RegistryMemoryDB struct {
	entries map[uuid.UUID]model.RegistryEntry
}

func CreateRegistryMemoryDB() (*RegistryMemoryDB, error) {
	return &RegistryMemoryDB{
		entries: make(map[uuid.UUID]model.RegistryEntry),
	}, nil
}

func (db *RegistryMemoryDB) CreateRegistryEntry(email string) (uuid.UUID, error) {
	id := uuid.New()
	db.entries[id] = model.RegistryEntry{
		Id:    id,
		Email: email,
	}
	return id, nil
}

func (db *RegistryMemoryDB) CheckIfRegistryEntryExistsByEmail(email string) (bool, error) {
	for _, entry := range db.entries {
		if entry.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (db *RegistryMemoryDB) CheckIfRegistryEntryExists(id uuid.UUID) (bool, error) {
	_, exists := db.entries[id]
	return exists, nil
}

func (db *RegistryMemoryDB) GetRegistryEntry(id uuid.UUID) (model.RegistryEntry, error) {
	entry, exists := db.entries[id]
	if !exists {
		return model.RegistryEntry{}, database.ErrKeyNotFound
	}
	return entry, nil
}

func (db *RegistryMemoryDB) GetRegistryEntryByEmail(email string) (model.RegistryEntry, error) {
	for _, entry := range db.entries {
		if entry.Email == email {
			return entry, nil
		}
	}
	return model.RegistryEntry{}, database.ErrKeyNotFound
}

func (db *RegistryMemoryDB) AddPersonalInfoToRegistryEntry(id uuid.UUID, personalInfo model.UserPersonalInfoRecord) error {
	entry, exists := db.entries[id]
	if !exists {
		return database.ErrKeyNotFound
	}
	entry.PersonalInfo = personalInfo
	db.entries[id] = entry
	return nil
}

func (db *RegistryMemoryDB) AddInterestsToRegistryEntry(id uuid.UUID, interests []string) error {
	entry, exists := db.entries[id]
	if !exists {
		return database.ErrKeyNotFound
	}
	entry.Interests = interests
	db.entries[id] = entry
	slog.Info("interests added to registry entry", slog.Any("rinterests", interests))
	return nil
}

func (db *RegistryMemoryDB) VerifyEmail(id uuid.UUID) error {
	entry, exists := db.entries[id]
	if !exists {
		return database.ErrKeyNotFound
	}
	entry.EmailVerified = true
	db.entries[id] = entry
	return nil
}

func (db *RegistryMemoryDB) DeleteRegistryEntry(id uuid.UUID) error {
	if _, exists := db.entries[id]; !exists {
		return database.ErrKeyNotFound
	}
	delete(db.entries, id)
	return nil
}
