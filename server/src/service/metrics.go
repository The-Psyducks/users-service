package service

import (
	// "fmt"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/model"
)

func (u *User) GetRegistrationMetrics(isAdmin bool) (*model.RegistrationSummaryMetrics, error) {
	if !isAdmin {
		return nil, app_errors.NewAppError(http.StatusForbidden, UserIsNotAdmin, ErrUserIsNotAdmin)
	}

	metris, err := u.registryDb.GetRegistrySummaryMetrics()
	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, err)
	}
	return metris, nil
}

func (u *User) GetLoginMetrics(isAdmin bool) (*model.LoginSummaryMetrics, error) {
	if !isAdmin {
		return nil, app_errors.NewAppError(http.StatusForbidden, UserIsNotAdmin, ErrUserIsNotAdmin)
	}

	metrics, err := u.userDb.GetLoginSummaryMetrics()
	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, err)
	}
	return metrics, nil
}

func (u *User) GetLocationMetrics(isAdmin bool) (*model.LocationMetrics, error) {
	if !isAdmin {
		return nil, app_errors.NewAppError(http.StatusForbidden, UserIsNotAdmin, ErrUserIsNotAdmin)
	}

	metrics, err := u.userDb.GetLocationMetrics()
	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, err)
	}
	return metrics, nil
}
