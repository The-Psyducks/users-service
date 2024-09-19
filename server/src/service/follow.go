package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/database"

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

	followerUUID, err := uuid.Parse(followerId)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error parsing followerId: %w", err))
	}

	err = u.userDb.FollowUser(followerUUID, userRecord.Id)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error following user: %w", err))
	}

	slog.Info("user followed succesfully", slog.String("followerId", followerId), slog.String("followingId", userRecord.Id.String()))
	return nil
}