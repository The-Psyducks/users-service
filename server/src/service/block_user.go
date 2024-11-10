package service

import (
	"fmt"
	"net/http"
	"users-service/src/app_errors"
	"github.com/google/uuid"
)

func (u *User) BlockUser(userId uuid.UUID, userSessionIsAdmin bool) error {
	if !userSessionIsAdmin {
		err := app_errors.NewAppError(http.StatusForbidden, UserIsNotAdmin, ErrUserIsNotAdmin)
		return err
	}

	if err := u.userDb.BlockUser(userId); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error blocking user: %w", err))
	}
	return nil
}

func (u *User) UnblockUser(userId uuid.UUID, userSessionIsAdmin bool) error {
	if !userSessionIsAdmin {
		err := app_errors.NewAppError(http.StatusForbidden, UserIsNotAdmin, ErrUserIsNotAdmin)
		return err
	}

	if err := u.userDb.UnblockUser(userId); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error unblocking user: %w", err))
	}
	return nil
}

func (u *User) CheckIfUserIsBlocked(userId uuid.UUID) (bool, error) {
	isBlocked, err := u.userDb.CheckIfUserIsBlocked(userId)
	if err != nil {
		return false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if user is blocked: %w", err))
	}
	return isBlocked, nil
}