package users_db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"users-service/src/constants"
	"users-service/src/database"
	"users-service/src/model"
)

type UsersPostgresDB struct {
	db *sqlx.DB
}

func CreateUsersPostgresDB(db *sqlx.DB) (*UsersPostgresDB, error) {
	dropDatabase := fmt.Sprintf(`
		DROP TABLE IF EXISTS %s CASCADE;
		DROP TABLE IF EXISTS %s CASCADE;
		`, "users", "user_interests")
	if _, err := db.Exec(dropDatabase); err != nil {
		return nil, fmt.Errorf("failed to drop database: %w", err)
	}

	schema := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(%d) NOT NULL UNIQUE,
			first_name VARCHAR(%d) NOT NULL,
			last_name VARCHAR(%d) NOT NULL,
			email VARCHAR(%d) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			location VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);

		CREATE TABLE IF NOT EXISTS user_interests (
		user_id UUID NOT NULL,
		interest VARCHAR(255) NOT NULL,
		PRIMARY KEY (user_id, interest),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`, constants.MaxUsernameLength, constants.MaxFirstNameLength, constants.MaxLastNameLength, constants.MaxEmailLength)
	

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	postgresDB := UsersPostgresDB{db}

	return &postgresDB, nil
}


func (postDB *UsersPostgresDB) CreateUser(data model.UserRecord) (model.UserRecord, error) {
	var user model.UserRecord
    query := `
        INSERT INTO users (username, first_name, last_name, email, password, location)
        VALUES (:username, :first_name, :last_name, :email, :password, :location)
        RETURNING id, username, first_name, last_name, email, password, location, created_at
    `

	rows, err := postDB.db.NamedQuery(query, data)
	if err != nil {
		return model.UserRecord{}, err
	}
	defer rows.Close()

    if rows.Next() {
        if err := rows.StructScan(&user); err != nil {
            return model.UserRecord{}, fmt.Errorf("error scanning user data: %w", err)
        }
    } else {
        return model.UserRecord{}, fmt.Errorf("error: no user created")
    }

	return user, nil
}

func (postDB *UsersPostgresDB) GetUserById(id uuid.UUID) (model.UserRecord, error) {
	var user model.UserRecord
	query := `SELECT * FROM users WHERE id = $1`
	err := postDB.db.Get(&user, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserRecord{}, database.ErrKeyNotFound
		}
		return model.UserRecord{}, fmt.Errorf("error fetching user by Id: %w", err)
	}
	return user, nil
}

func (postDB *UsersPostgresDB) GetUserByUsername(username string) (model.UserRecord, error) {
	var user model.UserRecord
	query := `SELECT * FROM users WHERE username = $1 LIMIT 1`
	err := postDB.db.Get(&user, query, username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserRecord{}, database.ErrKeyNotFound
		}
		return model.UserRecord{}, fmt.Errorf("error fetching user by username: %w", err)
	}
	return user, nil
}

func (postDB *UsersPostgresDB) GetUserByEmail(email string) (model.UserRecord, error) {
	var user model.UserRecord
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`
	err := postDB.db.Get(&user, query, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserRecord{}, database.ErrKeyNotFound
		}
		return model.UserRecord{}, fmt.Errorf("error fetching user by email: %w", err)
	}
	return user, nil
}

func (postDB *UsersPostgresDB) CheckIfUsernameExists(username string) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`
    err := postDB.db.QueryRow(query, username).Scan(&exists)

    if err != nil {
        return false, fmt.Errorf("error checking username existence: %w", err)
    }

    return exists, nil
}

func (postDB *UsersPostgresDB) CheckIfEmailExists(email string) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
    err := postDB.db.QueryRow(query, email).Scan(&exists)

    if err != nil {
        return false, fmt.Errorf("error checking email existence: %w", err)
    }

    return exists, nil
}

func (postDB *UsersPostgresDB) AssociateInterestsToUser(userId uuid.UUID, interests []string) error {
	var interestRecord model.Interest
	query := `
		INSERT INTO user_interests (user_id, interest)
		VALUES ($1, $2)
	`

	for _, interest := range interests {
		err := postDB.db.QueryRow(query, userId, interest).Scan(&interestRecord.Id, &interestRecord.Interest)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return fmt.Errorf("error inserting interest record: %w", err)
		}
	}
	return nil
}

func (postDB *UsersPostgresDB) GetInterestsForUserId(id uuid.UUID) ([]string, error) {
	var interests []string
	query := `
		SELECT interest
		FROM user_interests
		WHERE user_id = $1
	`

	rows, err := postDB.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error getting interests for user: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var interest string
		err := rows.Scan(&interest)
		if err != nil {
			return nil, fmt.Errorf("error scanning interest: %w", err)
		}
		interests = append(interests, interest)
	}

	return interests, nil
}