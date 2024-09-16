package interests_db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/google/uuid"
	"users-service/src/model"
)

type InterestsPostgresDB struct {
	db *sqlx.DB
}

func CreateInterestsPostgresDB(db *sqlx.DB) (*InterestsPostgresDB, error) {
	dropDatabase := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", "interests")
	if _, err := db.Exec(dropDatabase); err != nil {
		return nil, fmt.Errorf("failed to drop database: %w", err)
	}

	schema := `
    CREATE TABLE IF NOT EXISTS user_interests (
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

func (postDB *InterestsPostgresDB) AssociateInterestsToUser(userId uuid.UUID, interests []string) error {
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

func (postDB *InterestsPostgresDB) GetInterestsForUserId(id uuid.UUID) ([]string, error) {
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