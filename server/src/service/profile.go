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
)

func (u *User) GetUserProfile(userSessionId string, username string) (model.UserProfileResponse, error) {
	userRecord, err := u.userDb.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return model.UserProfileResponse{}, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return model.UserProfileResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if strings.EqualFold(userSessionId, userRecord.Id.String()) {
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
	interests, err := u.userDb.GetInterestsForUserId(user.Id)
	if err != nil {
		return model.UserProfileResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting interests from user: %w", err))
	}
	
	privateProfile, err := u.createUserPrivateProfileFromUserRecordAndInterests(user, interests)
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
