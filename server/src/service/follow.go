package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"users-service/src/app_errors"
	"users-service/src/database"
	"users-service/src/model"

	"github.com/google/uuid"
)

func sendNewFollowerNotification(followerId uuid.UUID, followingId uuid.UUID, token string) error {
	type NewFollowerNotification struct {
		UserId      uuid.UUID `json:"user_id"`
		FollowerId	uuid.UUID `json:"follower_id"`
	}

	url := "http://" + os.Getenv("NOTIF_HOST") + "/notification/followers-milestone"
	marshalledData, _ := json.Marshal(NewFollowerNotification{followingId, followerId})

	req, err := http.NewRequest("POST", url, bytes.NewReader(marshalledData))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	if err != nil {
		return errors.New("error creating request")
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return errors.New("error sending request, " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("error sending request, status code: "  + strconv.Itoa(resp.StatusCode))
	}

	slog.Info("Notification sent to ", followingId.String())
	return nil
}

func (u *User) FollowUser(followerId uuid.UUID, followingId uuid.UUID, token string) error {
	userRecord, err := u.userDb.GetUserById(followingId)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if followerId == userRecord.Id {
		return app_errors.NewAppError(http.StatusBadRequest, CantFollowYourself, fmt.Errorf("you can not following yourself"))
	}

	err = u.userDb.FollowUser(followerId, userRecord.Id)
	if err != nil {
		if errors.Is(err, database.ErrKeyAlreadyExists) {
			return app_errors.NewAppError(http.StatusBadRequest, AlreadyFollowing, err)
		}
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error following user: %w", err))
	}

	err = sendNewFollowerNotification(followerId, followingId, token)
	if err != nil {
		slog.Warn("Error sending notification to ", followingId.String(), slog.String("error: ", err.Error()))
	}

	slog.Info("user followed succesfully", slog.String("followerId", followerId.String()), slog.String("followingId", userRecord.Id.String()))
	return nil
}

func (u *User) UnfollowUser(followerId uuid.UUID, followingId uuid.UUID) error {
	userRecord, err := u.userDb.GetUserById(followingId)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	err = u.userDb.UnfollowUser(followerId, userRecord.Id)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return app_errors.NewAppError(http.StatusBadRequest, NotFollowing, err)
		}
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error unfollowing user: %w", err))
	}

	slog.Info("user unfollowed succesfully", slog.String("followerId", followerId.String()), slog.String("followingId", userRecord.Id.String()))
	return nil
}

// GetFollowers returns the followers of a user and if there are more to fetch
func (u *User) GetFollowers(id uuid.UUID, userSessionId uuid.UUID, timestamp string, skip, limit int) ([]model.UserProfileResponse, bool, error) {
	userRequested, err := u.userDb.GetUserById(id)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return nil, false, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return nil, false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if userRequested.Id != userSessionId {
		follows, err := u.userDb.CheckIfUserFollows(userSessionId, userRequested.Id)
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

	profiles, err := u.getUserProfilesFromUserRecords(followers, userSessionId)
	if err != nil {
		return nil, false, err
	}

	return profiles, hasMore, nil
}

// GetFollowers returns the user's a user is following and if there are more to fetch
func (u *User) GetFollowing(id uuid.UUID, userSessionId uuid.UUID, timestamp string, skip, limit int) ([]model.UserProfileResponse, bool, error) {
	userRecord, err := u.userDb.GetUserById(id)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return nil, false, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return nil, false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if userRecord.Id != userSessionId {
		follows, err := u.userDb.CheckIfUserFollows(userSessionId, userRecord.Id)
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

	profiles, err := u.getUserProfilesFromUserRecords(following, userSessionId)
	if err != nil {
		return nil, false, err
	}

	return profiles, hasMore, nil
}
