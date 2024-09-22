package service

import (
	"fmt"
	"log/slog"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/constants"
	"users-service/src/model"
)

func (u *User) checkIfEmailHasAccount(email string) (bool, error) {
	slog.Info("checking if email has account")

	user, err := u.userDb.CheckIfEmailExists(email)

	if err != nil {
		return false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if email exists: %w", err))
	}

	return user, nil

}

func (u *User) createNewRegistry(email string) (model.ResolveResponse, error) {
	registryId, err := u.registryDb.CreateRegistryEntry(email)
	if err != nil {
		return model.ResolveResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error creating registry entry: %w", err))
	}

	return model.ResolveResponse{
		NextAuthStep: constants.SignUpStep,
		Metadata: map[string]interface{}{
			"onboarding_step": constants.EmailVerificationStep,
			"registration_id": registryId.String(),
		},
	}, nil
}

func (u *User) resolveExistingRegistry(email string) (model.ResolveResponse, error) {
	registry, err := u.registryDb.GetRegistryEntryByEmail(email)
	if err != nil {
		return model.ResolveResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting registry entry: %w", err))
	}

	slog.Info("user email resolved successfully: it has registry", slog.String("email", email))
	return model.ResolveResponse{
		NextAuthStep: constants.SignUpStep,
		Metadata: map[string]interface{}{
			"onboarding_step": getStepForRegistryEntry(registry),
			"registration_id": registry.Id.String(),
		},
	}, nil
}

func (u *User) ResolveUserEmail(data model.ResolveRequest) (model.ResolveResponse, error) {
	slog.Info("resolving user email")

	if valErrs, err := u.userValidator.ValidateEmail(data.Email); err != nil {
		return model.ResolveResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error validating mail: %w", err))
	} else if len(valErrs) > 0 {
		return model.ResolveResponse{}, app_errors.NewAppValidationError(valErrs)
	}

	// chequeo de provider y verificacion del token
	


	hasAccount, err := u.checkIfEmailHasAccount(data.Email)
	if err != nil {
		return model.ResolveResponse{}, err
	}

	if hasAccount {
		slog.Info("user email resolved successfully: it has account", slog.String("email", data.Email))
		return model.ResolveResponse{
			NextAuthStep: constants.LoginStep,
			Metadata:     nil,
		}, nil
	}

	exists, err := u.registryDb.CheckIfRegistryEntryExistsByEmail(data.Email)
	if err != nil {
		return model.ResolveResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if registry entry exists: %w", err))
	}

	if exists {
		slog.Info("user email resolved successfully: it has registry entry", slog.String("email", data.Email))
		return u.resolveExistingRegistry(data.Email)
	}

	slog.Info("user email resolved successfully: it doesnt have account", slog.String("email", data.Email))
	return u.createNewRegistry(data.Email)
}