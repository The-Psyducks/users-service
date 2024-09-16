package service


import (
	"fmt"
	"errors"
	"net/http"
	"log/slog"
	"users-service/src/model"
	"users-service/src/database"
	"users-service/src/app_errors"
)


func (u *User) GetUserByUsername(username string) (model.UserResponse, error) {
	userRecord, err := u.userDb.GetUserByUsername(username)

	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return model.UserResponse{}, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	interests, err := u.interestDb.GetInterestsForUserId(userRecord.Id)
	if err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting interests from user: %w", err))
	}

	slog.Info("user retrieved succesfully", slog.String("username", userRecord.UserName))
	return createUserResponseFromUserRecordAndInterests(userRecord, interests), nil
}
