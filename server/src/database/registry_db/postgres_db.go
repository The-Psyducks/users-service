package registry_db

import (
	"database/sql"
	"fmt"
	"users-service/src/constants"
	"users-service/src/database"
	"users-service/src/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type RegistryPostgresDB struct {
	db *sqlx.DB
}

func CreateRegistryPostgresDB(db *sqlx.DB, test bool) (*RegistryPostgresDB, error) {
	if test {
		dropTables := `
			DROP TABLE IF EXISTS registry_interests CASCADE;
			DROP TABLE IF EXISTS registry_entries CASCADE;
		`
		if _, err := db.Exec(dropTables); err != nil {
			return nil, fmt.Errorf("failed to drop tables: %w", err)
		}
	}

	schema := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS registry_entries (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        email VARCHAR(%d) NOT NULL UNIQUE,
		identity_provider VARCHAR(255) DEFAULT NULL,
		email_verification_pin VARCHAR(10),
        email_verified BOOLEAN NOT NULL DEFAULT FALSE,
        first_name VARCHAR(%d) DEFAULT '',
        last_name VARCHAR(%d) DEFAULT '',
        username VARCHAR(%d) DEFAULT '',
        password TEXT DEFAULT '',
        location VARCHAR(255) DEFAULT '',
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		deleted_at TIMESTAMPTZ
    );

    CREATE TABLE IF NOT EXISTS registry_interests (
        registry_id UUID,
        interest VARCHAR(255),
        PRIMARY KEY (registry_id, interest),
        FOREIGN KEY (registry_id) REFERENCES registry_entries(id) ON DELETE CASCADE
    );`, constants.MaxEmailLength, constants.MaxFirstNameLength, constants.MaxLastNameLength, constants.MaxUsernameLength)

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &RegistryPostgresDB{db}, nil
}

func (db *RegistryPostgresDB) CreateRegistryEntry(email string, identityProvider *string) (uuid.UUID, error) {
	var id uuid.UUID
	err := db.db.QueryRow("INSERT INTO registry_entries (email, identity_provider) VALUES ($1, $2) RETURNING id", email, identityProvider).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create registry entry: %w", err)
	}
	return id, nil
}

func (db *RegistryPostgresDB) CheckIfRegistryEntryExistsByEmail(email string) (bool, error) {
	var exists bool
	err := db.db.QueryRow("SELECT EXISTS(SELECT 1 FROM registry_entries WHERE email = $1 AND deleted_at IS NULL)", email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if registry entry exists: %w", err)
	}
	return exists, nil
}

func (db *RegistryPostgresDB) CheckIfRegistryEntryExists(id uuid.UUID) (bool, error) {
	var exists bool
	err := db.db.QueryRow("SELECT EXISTS(SELECT 1 FROM registry_entries WHERE id = $1 AND deleted_at IS NULL)", id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if registry entry exists: %w", err)
	}
	return exists, nil
}

func (db *RegistryPostgresDB) GetRegistryEntry(id uuid.UUID) (model.RegistryEntry, error) {
	var entry model.RegistryEntry
	var personalInfo model.UserPersonalInfoRecord

	err := db.db.QueryRow(`
        SELECT id, email, email_verified, first_name, last_name, username, password, location 
        FROM registry_entries 
        WHERE id = $1`, id).Scan(
		&entry.Id, &entry.Email, &entry.EmailVerified,
		&personalInfo.FirstName, &personalInfo.LastName,
		&personalInfo.UserName, &personalInfo.Password,
		&personalInfo.Location)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.RegistryEntry{}, database.ErrKeyNotFound
		}
		return model.RegistryEntry{}, fmt.Errorf("failed to get registry entry: %w", err)
	}

	entry.PersonalInfo = personalInfo

	interests, err := db.getInterests(id)
	if err != nil {
		return model.RegistryEntry{}, fmt.Errorf("failed to get interests: %w", err)
	}
	entry.Interests = interests

	return entry, nil
}


func (db *RegistryPostgresDB) GetRegistryEntryByEmail(email string) (model.RegistryEntry, error) {
	var entry model.RegistryEntry
	var personalInfo model.UserPersonalInfoRecord

	err := db.db.QueryRow(`
        SELECT id, email, email_verified, first_name, last_name, username, password, location 
        FROM registry_entries 
        WHERE email = $1`, email).Scan(
		&entry.Id, &entry.Email, &entry.EmailVerified,
		&personalInfo.FirstName, &personalInfo.LastName,
		&personalInfo.UserName, &personalInfo.Password,
		&personalInfo.Location)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.RegistryEntry{}, database.ErrKeyNotFound
		}
		return model.RegistryEntry{}, fmt.Errorf("failed to get registry entry: %w", err)
	}

	entry.PersonalInfo = personalInfo

	interests, err := db.getInterests(entry.Id)
	if err != nil {
		return model.RegistryEntry{}, fmt.Errorf("failed to get interests: %w", err)
	}
	entry.Interests = interests

	return entry, nil
}

func (db *RegistryPostgresDB) AddPersonalInfoToRegistryEntry(id uuid.UUID, personalInfo model.UserPersonalInfoRecord) error {
	_, err := db.db.Exec(`
        UPDATE registry_entries 
        SET first_name = $2, last_name = $3, username = $4, password = $5, location = $6
        WHERE id = $1`,
		id, personalInfo.FirstName, personalInfo.LastName,
		personalInfo.UserName, personalInfo.Password, personalInfo.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return database.ErrKeyNotFound
		}
		return fmt.Errorf("failed to add personal info: %w", err)
	}
	return nil
}

func (db *RegistryPostgresDB) AddInterestsToRegistryEntry(id uuid.UUID, interests []string) error {
    var exists bool
    err := db.db.QueryRow("SELECT EXISTS(SELECT 1 FROM registry_interests WHERE registry_id = $1)", id).Scan(&exists)
    if err != nil {
        return fmt.Errorf("failed to check existing interests: %w", err)
    }

    if exists {
        return fmt.Errorf("interests already exist for registry entry with id %s", id)
    }

    for _, interest := range interests {
        _, err = db.db.Exec("INSERT INTO registry_interests (registry_id, interest) VALUES ($1, $2)", id, interest)
        if err != nil {
            return fmt.Errorf("failed to insert interest '%s': %w", interest, err)
        }
    }

    return nil
}

func (db *RegistryPostgresDB) SetEmailVerificationPin(id uuid.UUID, code string) error {
	res, err := db.db.Exec("UPDATE registry_entries SET email_verification_pin = $2 WHERE id = $1", id, code)
	if err != nil {
		return fmt.Errorf("failed to set email verification pin: %w", err)
	}
	if affected, err := res.RowsAffected(); err != nil || affected == 0 {
		return database.ErrKeyNotFound
	}
	return nil
}

func (db *RegistryPostgresDB) GetEmailVerificationPin(id uuid.UUID) (string, error) {
	var pin string
	err := db.db.QueryRow("SELECT email_verification_pin FROM registry_entries WHERE id = $1", id).Scan(&pin)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", database.ErrKeyNotFound
		}
		return "", fmt.Errorf("failed to get email verification pin: %w", err)
	}
	return pin, nil
}


func (db *RegistryPostgresDB) VerifyEmail(id uuid.UUID) error {
	_, err := db.db.Exec("UPDATE registry_entries SET email_verified = true WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}

func (db *RegistryPostgresDB) DeleteRegistryEntry(id uuid.UUID) error {
	_, err := db.db.Exec("UPDATE registry_entries SET deleted_at = now() WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to mark registry entry as deleted: %w", err)
	}
	return nil
}

func (db *RegistryPostgresDB) getInterests(id uuid.UUID) ([]string, error) {
	rows, err := db.db.Query("SELECT interest FROM registry_interests WHERE registry_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get interests: %w", err)
	}
	defer rows.Close()

	var interests []string
	for rows.Next() {
		var interest string
		if err := rows.Scan(&interest); err != nil {
			return nil, fmt.Errorf("failed to scan interest: %w", err)
		}
		interests = append(interests, interest)
	}

	return interests, nil
}
