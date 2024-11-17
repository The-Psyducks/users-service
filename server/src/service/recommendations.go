package service

import (
	"fmt"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/model"

	"github.com/google/uuid"
)

func (u *User) RecommendUsers(userSessionId uuid.UUID, timestamp string, skip int, limit int) ([]model.UserProfileResponse, bool, error) {
	users, hasMore, err := u.userDb.GetRecommendations(userSessionId, timestamp, skip, limit)
	if err != nil {
		err = app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting recommendations: %w", err))
		return nil, false, err
	}

	profiles, err := u.getUserProfilesFromUserRecords(users, userSessionId)
	if err != nil {
		return nil, false, err
	}

	return profiles, hasMore, nil
}