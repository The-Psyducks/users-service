package interests_db

import (
	// "database/sql"
	// "errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/google/uuid"
	// "users-service/src/constants"
	// "users-service/src/database"
	"users-service/src/model"
)

type InterestsPostgresDB struct {
	db *sqlx.DB
}

func CreateInterestsPostgresDB(databaseHost string, databasePort string, databaseName string, databasePassword string, databaseUser string) (*InterestsPostgresDB, error) {
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

	// dropDatabase := fmt.Sprintf("DROP TABLE IF EXISTS %s;", "interests")
	// if _, err := db.Exec(dropDatabase); err != nil {
	// 	return nil, fmt.Errorf("failed to drop database: %w", err)
	// }

	schema := `
    CREATE TABLE IF NOT EXISTS interests (
        user_id UUID NOT NULL,
        interest VARCHAR(255) NOT NULL,
		PRIMARY KEY (user_id, interest),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	postgresDB := InterestsPostgresDB{db}

	return &postgresDB, nil
}

func (postDB *InterestsPostgresDB) AssociateInterestsToUser(userId uuid.UUID, interests []int) ([]model.Interest, error) {
	var interestsRecords []model.Interest
	var interestRecord model.Interest
	query := `
		INSERT INTO user_interests (user_id, interest_id)
		VALUES ($1, $2)
		RETURNING user_id, interest_id
	`

	for _, interest := range interests {
		err := postDB.db.QueryRow(query, userId, interest).Scan(&interestRecord.Id, &interestRecord.Interest)
		if err != nil {
			return interestsRecords, fmt.Errorf("error inserting interest record: %w", err)
		}
		interestsRecords = append(interestsRecords, interestRecord)
	}

	return interestsRecords, nil
}
