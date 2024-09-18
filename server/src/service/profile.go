package service

import (
	"fmt"
	"errors"
	"strings"
	"log/slog"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/database"
	"users-service/src/model"
)

func (u *User) GetUserProfile(session_user_id string, username string) (model.UserPrivateProfile, error) {
	userRecord, err := u.userDb.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if strings.EqualFold(session_user_id, userRecord.Id.String()) {
		return u.getPublicProfile(userRecord), nil
	}
	return u.getPrivateProfile(userRecord)
}

func (u *User) getPrivateProfile(user model.UserRecord) (model.UserPrivateProfile, error) {
	slog.Info("user Private profile retrieved succesfully", slog.String("userId", user.Id.String()))
	interests, err := u.interestDb.GetInterestsForUserId(user.Id)
	if err != nil {
		return model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting interests from user: %w", err))
	}
	
	return createUserPrivateProfileFromUserRecordAndInterests(user, interests), nil
}

func (u *User) getPublicProfile(user model.UserRecord) model.UserPrivateProfile {
	slog.Info("user Public profile retrieved succesfully", slog.String("userId", user.Id.String()))
	return model.UserPrivateProfile{
		Id:       user.Id,
		FirstName: user.FirstName,
		LastName: user.LastName,
		UserName: user.UserName,
	}
}
