package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/auth"
	"users-service/src/database"
	"users-service/src/model"
)

func (u *User) loginValidUser(userRecord model.UserRecord, provider *string) (string, model.UserPrivateProfile, error) {
	authToken, err := auth.GenerateToken(userRecord.Id.String(), false)

	if err != nil {
		return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error generating token: %w", err))
	}

	privateProfile, err := u.createUserPrivateProfileFromUserRecord(userRecord)

	if err != nil {
		return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error generating private profile: %w", err))
	}

	if u.amqpQueue != nil {
		if err := u.sendLogInAttemptMessage(userRecord.Id.String(), true, provider); err != nil {
			slog.Warn("error publishing login attempt", slog.String("error", err.Error()))
		}
	}

	slog.Info("login information checked successfully", slog.String("username", userRecord.UserName))
	return authToken, privateProfile, nil	
}

func (u *User) LoginUser(data model.UserLoginRequest) (string, model.UserPrivateProfile, error) {
	slog.Info("checking login information")

	userRecord, err := u.userDb.GetUserByEmail(data.Email)

	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusNotFound, IncorrectUsernameOrPassword, errors.New("invalid username"))
		}
		return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if !checkPasswordHash(data.Password, userRecord.Password) {
		if err := u.sendLogInAttemptMessage(userRecord.Id.String(), false, nil); err != nil {
			slog.Warn("error publishing login attempt", slog.String("error", err.Error()))
		}
		return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusNotFound, IncorrectUsernameOrPassword, errors.New("invalid password"))
	}

	if isBlocked, err := u.CheckIfUserIsBlocked(userRecord.Id); err != nil {
		return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if user is blocked: %w", err))
	} else if isBlocked {
		return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusForbidden, UserBlocked, errors.New("user is blocked"))
	}

	return u.loginValidUser(userRecord, nil)
}

