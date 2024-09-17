package service

import (
	"fmt"
	"errors"
	"log/slog"
	"net/http"
	"users-service/src/auth"
	"users-service/src/model"
	"users-service/src/database"
	"users-service/src/app_errors"
)

func (u *User) CheckLoginCredentials(data model.UserLoginRequest) (string, error) {
	slog.Info("checking login information")

	userRecord, err := u.userDb.GetUserByUsername(data.UserName)

	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return "", app_errors.NewAppError(http.StatusNotFound, IncorrectUsernameOrPassword, errors.New("invalid username"))
		}
		return "", app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if !checkPasswordHash(data.Password, userRecord.Password) {
		return "", app_errors.NewAppError(http.StatusNotFound, IncorrectUsernameOrPassword, errors.New("invalid password"))
	}

	authToken, err := auth.GenerateToken(userRecord.UserName, true)

	if err != nil {
		return "", app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error generating token: %w", err))
	}

	slog.Info("login information checked successfully", slog.String("username", userRecord.UserName))
	return authToken, nil
}