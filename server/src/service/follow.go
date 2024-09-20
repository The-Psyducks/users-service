package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/database"
	"users-service/src/model"

	"github.com/google/uuid"
)

func (u *User) FollowUser(followerId string, followingUsername string) error {
	userRecord, err := u.userDb.GetUserByUsername(followingUsername)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if strings.EqualFold(followerId, userRecord.Id.String()) {
		return app_errors.NewAppError(http.StatusBadRequest, CantFollowYourself, fmt.Errorf("you can not following yourself"))
	}

	followerUUID, err := uuid.Parse(followerId)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error parsing followerId: %w", err))
	}

	err = u.userDb.FollowUser(followerUUID, userRecord.Id)
	if err != nil {
		if errors.Is(err, database.ErrKeyAlreadyExists) {
			return app_errors.NewAppError(http.StatusBadRequest, AlreadyFollowing, err)
		}
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error following user: %w", err))
	}

	slog.Info("user followed succesfully", slog.String("followerId", followerId), slog.String("followingId", userRecord.Id.String()))
	return nil
}

func (u *User) UnfollowUser(followerId string, followingUsername string) error {
	userRecord, err := u.userDb.GetUserByUsername(followingUsername)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	followerUUID, err := uuid.Parse(followerId)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error parsing followerId: %w", err))
	}

	err = u.userDb.UnfollowUser(followerUUID, userRecord.Id)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return app_errors.NewAppError(http.StatusBadRequest, NotFollowing, err)
		}
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error unfollowing user: %w", err))
	}

	slog.Info("user unfollowed succesfully", slog.String("followerId", followerId), slog.String("followingId", userRecord.Id.String()))
	return nil
}

func (u *User) GetFollowers(username string, userSessionId string) ([]model.FollowUserPublicProfile, error) {
	userRecord, err := u.userDb.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return nil, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return nil, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}
	fmt.Println("GetFollowers db")
	followers, err := u.userDb.GetFollowers(userRecord.Id)
	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting followers: %w", err))
	}
	
	fmt.Println("GetFollowersProfiles")
	profiles, err := u.getFollowersPublicProfilesFromUserRecords(followers, userSessionId)
	if err != nil {
		return nil, err
	}
	
	fmt.Println("Got theeeem")
	return profiles, nil
}

func (u *User) GetFollowing(username string, userSessionId string) ([]model.FollowUserPublicProfile, error) {
	userRecord, err := u.userDb.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return nil, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return nil, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	following, err := u.userDb.GetFollowing(userRecord.Id)
	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting following: %w", err))
	}

	profiles, err := u.getFollowersPublicProfilesFromUserRecords(following, userSessionId)
	if err != nil {
		return nil, err
	}
	
	return profiles, nil
}