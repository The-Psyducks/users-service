package service

import (
	"fmt"
	"users-service/src/model"

	"github.com/google/uuid"
)

// SearchUsers retrieves the users that have a username or name containing the text
// it sends first the ones that contain the text in the username
// then the ones that contain it in the name
// it also receives a timestamp, skip and limit to paginate the results
func (u *User) SearchUsers(userSessionId uuid.UUID, text string, timestamp string, skip int, limit int) ([]model.UserPublicProfileWithFollowStatus, bool, error) {
	users, hasMore, err := u.userDb.GetUsersWithUsernameContaining(text, timestamp, skip, limit)
	if err != nil {
		return nil, false, err
	}

	if len(users) < limit {
		var nameUsers []model.UserRecord
		remainingLimit := limit - len(users)
		remainingSkip := skip
		if len(users) == 0 { //case where I have to skip some users with name containing text 
			amntWithUsername, err := u.userDb.GetAmountOfUsersWithUsernameContaining(text)
			if err != nil {
				return nil, false, fmt.Errorf("error getting amount of users with username containing %s: %w", text, err)
			}
			remainingSkip = skip - amntWithUsername
		}

		nameUsers, hasMore, err = u.userDb.GetUsersWithNameContaining(text, timestamp, remainingSkip, remainingLimit)
		if err != nil {
			return nil, false, err
		}

		users = append(users, nameUsers...)
	}


	profiles, err := u.getFollowStatusPublicProfilesFromUserRecords(users, userSessionId)
	if err != nil {
		return nil, false, err
	}

	return profiles, hasMore, nil
}


