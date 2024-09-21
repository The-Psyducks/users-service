package users_db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	// "golang.org/x/crypto/bcrypt" //testing purposes

	"users-service/src/constants"
	"users-service/src/database"
	"users-service/src/model"
)

const (
	usersTable = "users"
	interestsTable = "user_interests"
	followersTable = "followers"
)

type UsersPostgresDB struct {
	db *sqlx.DB
}

func CreateUsersPostgresDB(db *sqlx.DB) (*UsersPostgresDB, error) {
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	postgresDB := UsersPostgresDB{db}

	// for testing purposes
	// postgresDB.createTestUsers()

	return &postgresDB, nil
}

func createTables(db *sqlx.DB) error {
	dropTables := fmt.Sprintf(`
		DROP TABLE IF EXISTS %s CASCADE;
		DROP TABLE IF EXISTS %s CASCADE;
		DROP TABLE IF EXISTS %s CASCADE;
		`, usersTable, interestsTable, followersTable)
	
	if _, err := db.Exec(dropTables); err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}
	
	schemaUsers := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
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
		`, usersTable, constants.MaxUsernameLength, constants.MaxFirstNameLength, constants.MaxLastNameLength, constants.MaxEmailLength)

	schemaInterests := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			user_id UUID NOT NULL,
			interest VARCHAR(255) NOT NULL,
			PRIMARY KEY (user_id, interest),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
			);
		`, interestsTable)
	
	schemaFollowers := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			follower_id UUID NOT NULL,
			following_id UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			PRIMARY KEY (follower_id, following_id),
			FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE
			);
		`, followersTable)

	if _, err := db.Exec(schemaUsers); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	if _, err := db.Exec(schemaInterests); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	if _, err := db.Exec(schemaFollowers); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
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

func (postDB *UsersPostgresDB) FollowUser(followerId uuid.UUID, followingId uuid.UUID) error {
	query := `
		INSERT INTO followers (follower_id, following_id)
		VALUES ($1, $2)
	`

	_, err := postDB.db.Exec(query, followerId, followingId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
            if pqErr.Code == "23505" { // Código de error para violación de unicidad en PostgreSQL
                return database.ErrKeyAlreadyExists
            }
        }
		return fmt.Errorf("error following user: %w", err)
	}
	return nil
}

func (postDB *UsersPostgresDB) UnfollowUser(followerId uuid.UUID, followingId uuid.UUID) error {
	query := `
		DELETE FROM followers
		WHERE follower_id = $1 AND following_id = $2
	`

	res, err := postDB.db.Exec(query, followerId, followingId)
	if err != nil {
		return fmt.Errorf("error unfollowing user: %w", err)
	}

   rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }

	if rowsAffected == 0 {
		return database.ErrKeyNotFound
	}
	return nil
}

func (postDB *UsersPostgresDB) CheckIfUserFollows(followerId string, followingId string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = $1 AND following_id = $2)`
	err := postDB.db.QueryRow(query, followerId, followingId).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("error checking if user follows: %w", err)
	}

	return exists, nil
}

func (postDB *UsersPostgresDB) GetAmountOfFollowers(userId uuid.UUID) (int, error) {
    var followers int
    query := `SELECT COUNT(*) FROM followers WHERE following_id = $1`
    err := postDB.db.Get(&followers, query, userId)

	if err != nil {
		return 0, fmt.Errorf("error getting amount of followers: %w", err)
	}
	return followers, nil
}


func (postDB *UsersPostgresDB) GetAmountOfFollowing(userId uuid.UUID) (int, error) {
	var following int
	query := `SELECT COUNT(*) FROM followers WHERE follower_id = $1`
	err := postDB.db.Get(&following, query, userId)

	if err != nil {
		return 0, fmt.Errorf("error getting amount of following: %w", err)
	}
	return following, nil
}

func (postDB *UsersPostgresDB) GetFollowers(userId uuid.UUID, timestamp string, skip int, limit int) ([]model.UserRecord, bool, error) {
	var followers []model.UserRecord
	query := `
		SELECT u.*
		FROM users u
		JOIN followers f ON u.id = f.follower_id
		WHERE f.following_id = $1
		AND f.created_at < $2
		ORDER BY f.created_at DESC
		OFFSET $3
		LIMIT $4
	`

	err := postDB.db.Select(&followers, query, userId, timestamp, skip, limit+1)
	if err != nil {
		return nil, false, fmt.Errorf("error getting followers: %w", err)
	}

	if len(followers) == limit+1 {
		return followers[:limit], true, nil
	}

	return followers, false, nil
}

func (postDB *UsersPostgresDB) GetFollowing(userId uuid.UUID, timestamp string, skip int, limit int) ([]model.UserRecord, bool, error) {
	var following []model.UserRecord
	query := `
		SELECT u.*
		FROM users u
		JOIN followers f ON u.id = f.following_id
		WHERE f.follower_id = $1
		AND f.created_at < $2
		ORDER BY f.created_at DESC
		OFFSET $3
		LIMIT $4
	`

	err := postDB.db.Select(&following, query, userId, timestamp, skip, limit+1)
	if err != nil {
		return nil, false, fmt.Errorf("error getting following: %w", err)
	}

	if len(following) == limit+1 {
		return following[:limit], true, nil
	}

	return following, false, nil
}

// // For testing purposes
// func hashPassword(password string) string {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		fmt.Println("error hashing password: ", err)
// 	}

// 	return string(hashedPassword)
// }

// func (postDB *UsersPostgresDB) createTestUsers() {
// 	users := []model.UserRecord{
// 		{
// 			UserName:  "Monke",
// 			FirstName: "Test",
// 			LastName:  "One",
// 			Email:     "monke@gmail.com",
// 			Password:  hashPassword("password"),
// 			Location:  "Test Location",
// 		},
// 		{
// 			UserName:  "Test",
// 			FirstName: "Test",
// 			LastName:  "Two",
// 			Email:     "test@gmail.com",
// 			Password:  hashPassword("password"),
// 			Location:  "Test Location",
// 		},
// 	}

// 	for _, user := range users {
// 		_, err := postDB.CreateUser(user)
// 		if err != nil {
// 			fmt.Println("error creating test user: ", err)
// 		}
// 	}		
// }