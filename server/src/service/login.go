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

	"github.com/google/uuid"
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

	if err := u.userDb.RegisterLoginAttempt(userRecord.Id, provider, true); err != nil {
		return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error registering login attempt: %w",
			err))
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
		if err := u.userDb.RegisterLoginAttempt(userRecord.Id, nil, false); err != nil {
			return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error registering login attempt: %w", err))
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

func (u *User) CheckIfUserIsBlocked(userId uuid.UUID) (bool, error) {
	isBlocked, err := u.userDb.CheckIfUserIsBlocked(userId)
	if err != nil {
		return false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if user is blocked: %w", err))
	}
	return isBlocked, nil
}
