package users_db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"users-service/src/constants"
	"users-service/src/database"
	"users-service/src/model"
)

type UsersPostgresDB struct {
	db *sqlx.DB
}

func CreateUsersPostgresDB(databaseHost string, databasePort string, databaseName string, databasePassword string, databaseUser string) (*UsersPostgresDB, error) {
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

	enableUUIDExtension := `
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    `
	if _, err := db.Exec(enableUUIDExtension); err != nil {
		return nil, fmt.Errorf("failed to enable uuid extension: %w", err)
	}

	dropDatabase := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", "users")
	if _, err := db.Exec(dropDatabase); err != nil {
		return nil, fmt.Errorf("failed to drop database: %w", err)
	}

	schema := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		username VARCHAR(%d) NOT NULL UNIQUE,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		email VARCHAR(%d) NOT NULL UNIQUE,
		password TEXT NOT NULL,
		location VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	);
	
	CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`, constants.MaxUsernameLength, constants.MaxEmailLength)
	

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

func (postDB *UsersPostgresDB) GetUserById(id string) (model.UserRecord, error) {
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

func (postDB *UsersPostgresDB) CheckIfUsernameExists(username string) (bool, error) {
    var count int
    query := `SELECT COUNT(*) FROM users WHERE LOWER(username) = LOWER($1)`
    err := postDB.db.QueryRow(query, username).Scan(&count)

    if err != nil {
        return false, fmt.Errorf("error checking username existence: %w", err)
    }

    return count > 0, nil
}

func (postDB *UsersPostgresDB) CheckIfMailExists(mail string) (bool, error) {
    var count int
    query := `SELECT COUNT(*) FROM users WHERE email = $1`
    err := postDB.db.QueryRow(query, mail).Scan(&count)

    if err != nil {
        return false, fmt.Errorf("error checking email existence: %w", err)
    }

    return count > 0, nil
}
