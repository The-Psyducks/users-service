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

func CreateRegistryPostgresDB(databaseHost, databasePort, databaseName, databasePassword, databaseUser string) (*RegistryPostgresDB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		databaseUser,
		databasePassword,
		databaseHost,
		databasePort,
		databaseName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	enableUUIDExtension := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	if _, err := db.Exec(enableUUIDExtension); err != nil {
		return nil, fmt.Errorf("failed to enable uuid extension: %w", err)
	}

	dropTables := `
    DROP TABLE IF EXISTS registry_interests CASCADE;
    DROP TABLE IF EXISTS registry_entries CASCADE;
	`
	if _, err := db.Exec(dropTables); err != nil {
		return nil, fmt.Errorf("failed to drop tables: %w", err)
	}

	schema := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS registry_entries (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        email VARCHAR(%d) NOT NULL UNIQUE,
        email_verified BOOLEAN NOT NULL DEFAULT FALSE,
        first_name VARCHAR(%d),
        last_name VARCHAR(%d),
        username VARCHAR(%d) UNIQUE,
        password TEXT,
        location VARCHAR(255)
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

func (db *RegistryPostgresDB) CreateRegistryEntry(email string) (uuid.UUID, error) {
	var id uuid.UUID
	err := db.db.QueryRow("INSERT INTO registry_entries (email) VALUES ($1) RETURNING id", email).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create registry entry: %w", err)
	}
	return id, nil
}

func (db *RegistryPostgresDB) CheckIfRegistryEntryExistsByEmail(email string) (bool, error) {
	var exists bool
	err := db.db.QueryRow("SELECT EXISTS(SELECT 1 FROM registry_entries WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if registry entry exists: %w", err)
	}
	return exists, nil
}

func (db *RegistryPostgresDB) CheckIfRegistryEntryExists(id uuid.UUID) (bool, error) {
	var exists bool
	err := db.db.QueryRow("SELECT EXISTS(SELECT 1 FROM registry_entries WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if registry entry exists: %w", err)
	}
	return exists, nil
}

func (db *RegistryPostgresDB) GetRegistryEntry(id uuid.UUID) (model.RegistryEntry, error) {
	var entry model.RegistryEntry
	var personalInfo struct {
		FirstName sql.NullString
		LastName  sql.NullString
		UserName  sql.NullString
		Password  sql.NullString
		Location  sql.NullString
	}

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

	entry.PersonalInfo = model.UserPersonalInfoRecord{
		FirstName: personalInfo.FirstName.String,
		LastName:  personalInfo.LastName.String,
		UserName:  personalInfo.UserName.String,
		Password:  personalInfo.Password.String,
		Location:  personalInfo.Location.String,
	}

	interests, err := db.getInterests(id)
	if err != nil {
		return model.RegistryEntry{}, fmt.Errorf("failed to get interests: %w", err)
	}
	entry.Interests = interests

	return entry, nil
}

func (db *RegistryPostgresDB) GetRegistryEntryByEmail(email string) (model.RegistryEntry, error) {
	var entry model.RegistryEntry
	var personalInfo struct {
		FirstName sql.NullString
		LastName  sql.NullString
		UserName  sql.NullString
		Password  sql.NullString
		Location  sql.NullString
	}

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

	entry.PersonalInfo = model.UserPersonalInfoRecord{
		FirstName: personalInfo.FirstName.String,
		LastName:  personalInfo.LastName.String,
		UserName:  personalInfo.UserName.String,
		Password:  personalInfo.Password.String,
		Location:  personalInfo.Location.String,
	}

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
    tx, err := db.db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    _, err = tx.Exec("DELETE FROM registry_interests WHERE registry_id = $1", id)
    if err != nil {
        return fmt.Errorf("failed to delete existing interests: %w", err)
    }

    stmt, err := tx.Prepare("INSERT INTO registry_interests (registry_id, interest) VALUES ($1, $2)")
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()

    for _, interest := range interests {
        _, err = stmt.Exec(id, interest)
        if err != nil {
            return fmt.Errorf("failed to insert interest '%s': %w", interest, err)
        }
    }

    if err = tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

func (db *RegistryPostgresDB) VerifyEmail(id uuid.UUID) error {
	_, err := db.db.Exec("UPDATE registry_entries SET email_verified = true WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}

func (db *RegistryPostgresDB) DeleteRegistryEntry(id uuid.UUID) error {
	_, err := db.db.Exec("DELETE FROM registry_entries WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete registry entry: %w", err)
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
