package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"users-service/src/app_errors"
	"users-service/src/auth"
	"users-service/src/constants"
	"users-service/src/database"
	"users-service/src/model"

	amqp "github.com/rabbitmq/amqp091-go"
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

	// if err := u.userDb.RegisterLoginAttempt(userRecord.Id, provider, true); err != nil {
	// 	return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error registering login attempt: %w",
	// 		err))
	// }
	if u.amqpQueue == nil {
		return authToken, privateProfile, nil
	}
	
	if provider == nil {
		fiero := constants.InternalProvider
		provider = &fiero
	}

	loginAttempt, err := json.Marshal(model.LoginAttempt{
										Succesfull: true,
										UserId:     userRecord.Id.String(),
										Provider:   *provider,
										Timestamp:  time.Now().GoString(),
									})
	if err != nil {
		return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error marshalling login attempt: %w", err))
	}

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        loginAttempt,
		DeliveryMode: amqp.Persistent,
	}

	err = u.amqpQueue.Publish("", os.Getenv("CLOUDAMQP_QUEUE"), false, false, message)

	if err != nil {
		slog.Error("error publishing login attempt", slog.String("error", err.Error()))
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
		// if err := u.userDb.RegisterLoginAttempt(userRecord.Id, nil, false); err != nil {
		// 	return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error registering login attempt: %w", err))
		// }
		return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusNotFound, IncorrectUsernameOrPassword, errors.New("invalid password"))
	}

	if isBlocked, err := u.CheckIfUserIsBlocked(userRecord.Id); err != nil {
		return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if user is blocked: %w", err))
	} else if isBlocked {
		return "", model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusForbidden, UserBlocked, errors.New("user is blocked"))
	}

	return u.loginValidUser(userRecord, nil)
}

