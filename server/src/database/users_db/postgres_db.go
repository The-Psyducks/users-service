package users_db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	// "golang.org/x/crypto/bcrypt" //testing purposes

	"users-service/src/constants"
	"users-service/src/database"
	"users-service/src/model"
)

const (
	usersTable     = "users"
	interestsTable = "user_interests"
	followersTable = "followers"
)

type UsersPostgresDB struct {
	db *sqlx.DB
}

func CreateUsersPostgresDB(db *sqlx.DB, test bool) (*UsersPostgresDB, error) {
	if err := createTables(db, test); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	postgresDB := UsersPostgresDB{db}

	// for testing purposes
	// postgresDB.createTestUsers()

	return &postgresDB, nil
}

func createTables(db *sqlx.DB, test bool) error {
	if test {
		dropTables := fmt.Sprintf(`
			DROP TABLE IF EXISTS %s CASCADE;
			DROP TABLE IF EXISTS %s CASCADE;
			DROP TABLE IF EXISTS %s CASCADE;
			`, usersTable, interestsTable, followersTable)

		if _, err := db.Exec(dropTables); err != nil {
			return fmt.Errorf("failed to drop database: %w", err)
		}
	}

	schemaUsers := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(%d) NOT NULL UNIQUE,
			picture_path TEXT DEFAULT '',
			first_name VARCHAR(%d) NOT NULL,
			last_name VARCHAR(%d) NOT NULL,
			email VARCHAR(%d) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			location VARCHAR(255) NOT NULL,
			blocked BOOLEAN DEFAULT FALSE,
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

func (postDB *UsersPostgresDB) associateInterestsToUser(userId uuid.UUID, interests []string) ([]string, error) {
	var insertedInterests []string
	query := `
		INSERT INTO user_interests (user_id, interest)
		VALUES ($1, $2)
		RETURNING interest;
	`

	for _, interest := range interests {
		var interestRecord string
		err := postDB.db.QueryRow(query, userId, interest).Scan(&interestRecord)
		if err != nil {
			return nil, fmt.Errorf("error inserting interest record: %w", err)
		}
		insertedInterests = append(insertedInterests, interestRecord)
	}

	return insertedInterests, nil
}


func (postDB *UsersPostgresDB) updateUserInterests(userId uuid.UUID, interests []string) ([]string, error) {
	query := `DELETE FROM user_interests WHERE user_id = $1`
	_, err := postDB.db.Exec(query, userId)
	if err != nil {
		return nil, fmt.Errorf("error deleting user interests: %w", err)
	}

	return postDB.associateInterestsToUser(userId, interests)
}

func (postDB *UsersPostgresDB) CreateUser(data model.UserRecord) (model.UserRecord, error) {
	var user model.UserRecord
	query := `
        INSERT INTO users (username, first_name, last_name, email, password, location)
        VALUES (:username, :first_name, :last_name, :email, :password, :location)
        RETURNING id, username, first_name, last_name, email, password, location, created_at;
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
	}

	user.Interests, err = postDB.associateInterestsToUser(user.Id, data.Interests)
	if err != nil {
		return model.UserRecord{}, fmt.Errorf("error associating interests to user: %w", err)
	}
	return user, nil
}

func (postDB *UsersPostgresDB) ModifyUser(id uuid.UUID, data model.UpdateUserPrivateProfile) (model.UserRecord, error) {
	var user model.UserRecord
	query := `
		UPDATE users
		SET username = :username, first_name = :first_name, last_name = :last_name, location = :location, picture_path = :picture_path
		WHERE id = :id
		RETURNING id, username, first_name, last_name, email, location, picture_path
	`

	rows, err := postDB.db.NamedQuery(query, map[string]interface{}{
		"id":         id,
		"username":   data.UserName,
		"first_name": data.FirstName,
		"last_name":  data.LastName,
		"location":   data.Location,
		"picture_path": data.PicturePath,
	})
	if err != nil {
		return model.UserRecord{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return model.UserRecord{}, fmt.Errorf("error scanning user data: %w", err)
		}
	} else {
		return model.UserRecord{}, fmt.Errorf("error: no user updated")
	}

	user.Interests, err = postDB.updateUserInterests(id, data.Interests)
	if err != nil {
		return model.UserRecord{}, fmt.Errorf("error updating user interests: %w", err)
	}
	return user, nil
}

// // For testing purposes
// func (postDB *UsersPostgresDB) PrintAllUsers() error {
//     var users []model.UserRecord
//     query := `SELECT * FROM users`

//     err := postDB.db.Select(&users, query)
//     if err != nil {
//         return fmt.Errorf("error fetching users: %w", err)
//     }

//     // Imprimir cada usuario
//     fmt.Println("All users in the database:")
//     for _, user := range users {
//         fmt.Printf("ID: %s, Username: %s, First Name: %s, Last Name: %s, Email: %s, Location: %s, Created At: %s\n",
//             user.Id, user.UserName, user.FirstName, user.LastName, user.Email, user.Location, user.CreatedAt)
//     }

//     return nil
// }

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

	user.Interests, err = postDB.getInterestsForUserId(id)
	if err != nil {
		return model.UserRecord{}, fmt.Errorf("error getting interests for user: %w", err)
	}
	fmt.Println("user in db: ", user)
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

	user.Interests, err = postDB.getInterestsForUserId(user.Id)
	if err != nil {
		return model.UserRecord{}, fmt.Errorf("error getting interests for user: %w", err)
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

func (postDB *UsersPostgresDB) getInterestsForUserId(id uuid.UUID) ([]string, error) {
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

func (postDB *UsersPostgresDB) CheckIfUserFollows(followerId uuid.UUID, followingId uuid.UUID) (bool, error) {
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

	for i := range followers {
		followers[i].Interests, err = postDB.getInterestsForUserId(followers[i].Id)
		if err != nil {
			return nil, false, fmt.Errorf("error getting interests for user: %w", err)
		}
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

func (postDB *UsersPostgresDB) GetAmountOfFollowersInTimeRange(userId uuid.UUID, startTime, endTime time.Time) (int, error) {
	var followers int
	query := `SELECT COUNT(*) FROM followers WHERE following_id = $1 AND created_at >= $2 AND created_at <= $3`
	err := postDB.db.Get(&followers, query, userId, startTime, endTime)

	if err != nil {
		return 0, fmt.Errorf("error getting amount of followers in time range: %w", err)
	}
	return followers, nil
}

func (postDB *UsersPostgresDB) GetAllUsers(timestamp string, skip int, limit int) ([]model.UserRecord, bool, error) {
	var users []model.UserRecord
	query := `
		SELECT *
		FROM users
		WHERE created_at < $1
		ORDER BY created_at DESC
		OFFSET $2
		LIMIT $3
	`

	err := postDB.db.Select(&users, query, timestamp, skip, limit+1)
	if err != nil {
		return nil, false, fmt.Errorf("error getting all users: %w", err)
	}

	if len(users) == limit+1 {
		return users[:limit], true, nil
	}

	for i := range users {
		users[i].Interests, err = postDB.getInterestsForUserId(users[i].Id)
		if err != nil {
			return nil, false, fmt.Errorf("error getting interests for user: %w", err)
		}
	}

	return users, false, nil
}

func (postDB *UsersPostgresDB) GetUsersWithUsernameContaining(text string, timestamp string, skip int, limit int) ([]model.UserRecord, bool, error) {
	var users []model.UserRecord
	query := `
		SELECT *
		FROM users
		WHERE username ILIKE $1
		AND created_at < $2
		ORDER BY created_at DESC
		OFFSET $3
		LIMIT $4
	`

	err := postDB.db.Select(&users, query, "%"+text+"%", timestamp, skip, limit+1)
	if err != nil {
		return nil, false, fmt.Errorf("error getting users with username containing: %w", err)
	}

	if len(users) == limit+1 {
		return users[:limit], true, nil
	}

	for i := range users {
		users[i].Interests, err = postDB.getInterestsForUserId(users[i].Id)
		if err != nil {
			return nil, false, fmt.Errorf("error getting interests for user: %w", err)
		}
	}

	return users, false, nil
}

func (postDB *UsersPostgresDB) GetAmountOfUsersWithUsernameContaining(text string) (int, error) {
	var amount int
	query := `SELECT COUNT(*) FROM users WHERE username ILIKE $1`
	err := postDB.db.Get(&amount, query, "%"+text+"%")

	if err != nil {
		return 0, fmt.Errorf("error getting amount of users with username containing: %w", err)
	}

	return amount, nil
}

func (postDB *UsersPostgresDB) GetUsersWithOnlyNameContaining(text string, timestamp string, skip int, limit int) ([]model.UserRecord, bool, error) {
	var users []model.UserRecord
	query := `
		SELECT *
		FROM users
		WHERE (first_name ILIKE $1 OR last_name ILIKE $1)
		AND username NOT ILIKE $1
		AND created_at < $2
		ORDER BY created_at DESC
		OFFSET $3
		LIMIT $4
	`

	err := postDB.db.Select(&users, query, "%"+text+"%", timestamp, skip, limit+1)
	if err != nil {
		return nil, false, fmt.Errorf("error getting users with name containing: %w", err)
	}

	if len(users) == limit+1 {
		return users[:limit], true, nil
	}

	return users, false, nil
}

func (postDB *UsersPostgresDB) GetRecommendations(userId uuid.UUID, timestamp string, skip int, limit int) ([]model.UserRecord, bool, error) {
	var users []model.UserRecord
	query := `
		WITH temp AS (
			SELECT DISTINCT ON (id)
				id, 
				username, 
				first_name, 
				last_name, 
				email, 
				location, 
				created_at, 
				priority
			FROM (
				(SELECT u.id, u.username, u.first_name, u.last_name, u.email, u.location, u.created_at, 1 AS priority
				FROM users u
				LEFT JOIN followers f ON u.id = f.following_id AND f.follower_id = $1
				JOIN user_interests ui ON u.id = ui.user_id
				JOIN users u2 ON u.location = u2.location AND u2.id = $1
				WHERE u.id != $1
				AND f.following_id IS NULL
				AND EXISTS (
					SELECT 1
					FROM user_interests ui2
					WHERE ui2.user_id = $1
						AND ui2.interest = ui.interest
				)
				AND u.created_at < $2
				ORDER BY u.created_at DESC
				LIMIT $4)

				UNION ALL

				(SELECT u.id, u.username, u.first_name, u.last_name, u.email, u.location, u.created_at, 2 AS priority
				FROM users u
				LEFT JOIN followers f ON u.id = f.following_id AND f.follower_id = $1
				JOIN users u2 ON u.location = u2.location AND u2.id = $1
				WHERE u.id != $1
				AND f.following_id IS NULL
				AND u.created_at < $2
				ORDER BY u.created_at DESC
				LIMIT $4)

				UNION ALL

				(SELECT u.id, u.username, u.first_name, u.last_name, u.email, u.location, u.created_at, 3 AS priority
				FROM users u
				LEFT JOIN followers f ON u.id = f.following_id AND f.follower_id = $1
				JOIN user_interests ui ON u.id = ui.user_id
				WHERE u.id != $1
				AND f.following_id IS NULL
				AND EXISTS (
					SELECT 1
					FROM user_interests ui2
					WHERE ui2.user_id = $1
						AND ui2.interest = ui.interest
				)
				AND u.created_at < $2
				ORDER BY u.created_at DESC
				LIMIT $4)
			)
		)
		SELECT 
			p.id, 
			p.username, 
			p.first_name, 
			p.last_name, 
			p.email, 
			p.location
		FROM temp p
		ORDER BY p.priority ASC, p.created_at DESC
		OFFSET $3
		LIMIT $4;
	`

	err := postDB.db.Select(&users, query, userId, timestamp, skip, limit+1)
	if err != nil {
		return nil, false, fmt.Errorf("error getting users with name containing: %w", err)
	}

	if len(users) == limit+1 {
		return users[:limit], true, nil
	}

	for i := range users {
		users[i].Interests, err = postDB.getInterestsForUserId(users[i].Id)
		if err != nil {
			return nil, false, fmt.Errorf("error getting interests for user: %w", err)
		}
	}

	return users, false, nil
}

func (postDB *UsersPostgresDB) BlockUser(userId uuid.UUID, reason string) error {
    query := `UPDATE users SET blocked = TRUE WHERE id = $1`
	_, err := postDB.db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error blocking user: %w", err)
	}
	return nil
}

func (postDB *UsersPostgresDB) UnblockUser(userId uuid.UUID) error {
    query := `UPDATE users SET blocked = FALSE WHERE id = $1`
	_, err := postDB.db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error unblocking user: %w", err)
	}
	return nil
}

func (postDB *UsersPostgresDB) CheckIfUserIsBlocked(userId uuid.UUID) (bool, error) {
    var count int
    query := `SELECT COUNT(*) FROM users WHERE id = $1 AND blocked = TRUE`
	err := postDB.db.Get(&count, query, userId)
	if err != nil {
		return false, fmt.Errorf("error checking if user is blocked: %w", err)
	}
	return count > 0, nil
}

// For testing purposes
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
