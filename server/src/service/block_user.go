package service

import (
	"fmt"
	"net/http"
	"users-service/src/app_errors"

	"github.com/google/uuid"
)

func (u *User) BlockUser(userId uuid.UUID, userSessionIsAdmin bool, reason string) error {
	if !userSessionIsAdmin {
		err := app_errors.NewAppError(http.StatusForbidden, UserIsNotAdmin, ErrUserIsNotAdmin)
		return err
	}

	if err := u.userDb.BlockUser(userId, reason); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error blocking user: %w", err))
	}

	if u.amqpQueue != nil {
		if err := u.sendUserBlockedMessage(userId.String(), reason); err != nil {
			return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error sending user blocked message: %w", err))
		}
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

	if u.amqpQueue != nil {
		if err := u.sendUserUnblockedMessage(userId.String()); err != nil {
			return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error sending user blocked message: %w", err))
		}
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
