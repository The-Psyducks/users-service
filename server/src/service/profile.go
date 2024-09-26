package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"users-service/src/app_errors"
	"users-service/src/database"
	"users-service/src/model"

	"github.com/google/uuid"
)

func (u *User) GetUserProfileById(userSessionId string, id uuid.UUID) (model.UserProfileResponse, error) {
	userRecord, err := u.userDb.GetUserById(id)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return model.UserProfileResponse{}, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return model.UserProfileResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if strings.EqualFold(userSessionId, id.String()) {
		return u.getPrivateProfile(userRecord)
	}
	return u.getPublicProfile(userRecord, userSessionId)
}

func (u *User) getAmountOfFollowersAndFollowing(user model.UserRecord) (int, int, error) {
	followers, err := u.userDb.GetAmountOfFollowers(user.Id)
	if err != nil {
		return 0, 0, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting amount of followers: %w", err))
	}

	following, err := u.userDb.GetAmountOfFollowing(user.Id)
	if err != nil {
		return 0, 0, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting amount of following: %w", err))
	}

	return followers, following, nil
}

func (u *User) getPrivateProfile(user model.UserRecord) (model.UserProfileResponse, error) {
	privateProfile, err := u.createUserPrivateProfileFromUserRecord(user)
	if err != nil {
		return model.UserProfileResponse{}, err
	}
	
	slog.Info("user Private profile retrieved succesfully", slog.String("userId", user.Id.String()))
	return model.UserProfileResponse{
		OwnProfile: true,
		Follows:    false,
		Profile:    privateProfile,
	}, nil
}

func (u *User) getPublicProfile(user model.UserRecord, session_user_id string) (model.UserProfileResponse, error) {
	profile, err := u.generateUserPublicProfileFromUserRecord(user)
	if err != nil {
		return model.UserProfileResponse{}, err
	}
	
	follows, err := u.userDb.CheckIfUserFollows(session_user_id, user.Id.String())
	if err != nil {
		return model.UserProfileResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if user follows: %w", err))
	}
	
	slog.Info("user Public profile retrieved succesfully", slog.String("userId", user.Id.String()))
	return model.UserProfileResponse{
		OwnProfile: false,
		Follows:    follows,
		Profile:    profile,
	}, nil
}
