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

// GetFollowers returns the followers of a user and if there are more to fetch
func (u *User) GetFollowers(username, userSessionId, timestamp string, skip, limit int) ([]model.FollowUserPublicProfile, bool, error) {
	userRequested, err := u.userDb.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return nil, false, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return nil, false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}
	
	if userRequested.Id.String() != userSessionId {
		follows, err := u.userDb.CheckIfUserFollows(userSessionId, userRequested.Id.String())
		if err != nil {
			return nil, false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if user follows: %w", err))
		}
		
		if !follows {
			return nil, false, app_errors.NewAppError(http.StatusForbidden, NotFollowing, fmt.Errorf("user does not follow the user"))
		}
	}
	
	followers, hasMore, err := u.userDb.GetFollowers(userRequested.Id, timestamp, skip, limit)
	if err != nil {
		return nil, false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting followers: %w", err))
	}
	
	profiles, err := u.getFollowersPublicProfilesFromUserRecords(followers, userSessionId)
	if err != nil {
		return nil, false, err
	}
	
	return profiles, hasMore, nil
}

// GetFollowers returns the user's a user is following and if there are more to fetch
func (u *User) GetFollowing(username, userSessionId, timestamp string, skip, limit int) ([]model.FollowUserPublicProfile, bool, error) {
	userRecord, err := u.userDb.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return nil, false, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return nil, false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if userRecord.Id.String() != userSessionId {
		follows, err := u.userDb.CheckIfUserFollows(userSessionId, userRecord.Id.String())
		if err != nil {
			return nil, false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if user follows: %w", err))
		}

		if !follows {
			return nil, false, app_errors.NewAppError(http.StatusForbidden, NotFollowing, fmt.Errorf("user does not follow the user"))
		}
	}

	following, hasMore, err := u.userDb.GetFollowing(userRecord.Id, timestamp, skip, limit)
	if err != nil {
		return nil, false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting following: %w", err))
	}

	profiles, err := u.getFollowersPublicProfilesFromUserRecords(following, userSessionId)
	if err != nil {
		return nil, false, err
	}
	
	return profiles, hasMore, nil
}