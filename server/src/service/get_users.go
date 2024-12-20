package service

import (
	"errors"
	"fmt"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/database"
	"users-service/src/model"

	"github.com/google/uuid"
)

// SearchUsers retrieves the users that have a username or name containing the text
// it sends first the ones that contain the text in the username
// then the ones that contain it in the name
// it also receives a timestamp, skip and limit to paginate the results
func (u *User) SearchUsers(userSessionId uuid.UUID, text string, timestamp string, skip int, limit int) ([]model.UserProfileResponse, bool, error) {
	users, hasMore, err := u.userDb.GetUsersWithUsernameContaining(text, timestamp, skip, limit)
	if err != nil {
		err = app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting users with username containing %s: %w", text, err))
		return nil, false, err
	}

	if len(users) < limit {
		var nameUsers []model.UserRecord
		remainingLimit := limit - len(users)
		remainingSkip := skip
		if len(users) == 0 { //case where I have to skip some users with name containing text
			amntWithUsername, err := u.userDb.GetAmountOfUsersWithUsernameContaining(text)
			if err != nil {
				err = app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting amount of users with username containing %s: %w", text, err))
				return nil, false, err
			}
			remainingSkip = skip - amntWithUsername
		}

		nameUsers, hasMore, err = u.userDb.GetUsersWithOnlyNameContaining(text, timestamp, remainingSkip, remainingLimit)
		if err != nil {
			err = app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting users with name containing %s: %w", text, err))
			return nil, false, err
		}

		users = append(users, nameUsers...)
	}

	profiles, err := u.getUserProfilesFromUserRecords(users, userSessionId)
	if err != nil {
		return nil, false, err
	}

	return profiles, hasMore, nil
}

// GetAllUsers retrieves all the users in the database, it is just for admins
// it also receives a timestamp, skip and limit to paginate the results
func (u *User) GetAllUsers(userSessionIsAdmin bool, timestamp string, skip int, limit int) ([]model.UserPublicProfile, bool, error) {
	if !userSessionIsAdmin {
		err := app_errors.NewAppError(http.StatusForbidden, UserIsNotAdmin, ErrUserIsNotAdmin)
		return nil, false, err
	}
	users, hasMore, err := u.userDb.GetAllUsers(timestamp, skip, limit)
	if err != nil {
		err = app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting all users: %w", err))
		return nil, false, err
	}

	profiles, err := u.getPublicProfilesFromUserRecords(users)
	if err != nil {
		err = app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting public profiles from user records: %w", err))
		return nil, false, err
	}

	return profiles, hasMore, nil
}

func (u *User) GetUserInformation(userSessionId uuid.UUID, userSessionIsAdmin bool, id uuid.UUID) (model.UserInformationResponse, error) {
	if !userSessionIsAdmin {
		return model.UserInformationResponse{}, app_errors.NewAppError(http.StatusForbidden, UserIsNotAdmin, ErrUserIsNotAdmin)
	}
	userRecord, err := u.userDb.GetUserById(id)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			err = app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		} else {
			err = app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
		}
		return model.UserInformationResponse{}, err
	}

	profile, err := u.getPrivateProfile(userRecord)
	if err != nil {
		return model.UserInformationResponse{}, err
	}

	isBlocked, err := u.userDb.CheckIfUserIsBlocked(id)
	if err != nil {
		err = app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if user is blocked: %w", err))
		return model.UserInformationResponse{}, err
	}

	return model.UserInformationResponse{
		IsBlocked: isBlocked,
		Profile:   profile.Profile,
	}, nil
}