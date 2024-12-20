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

func (u *User) createNewRegistry(email string, identityProvider *string) (model.ResolveResponse, error) {
	registryId, err := u.registryDb.CreateRegistryEntry(email, identityProvider)
	if err != nil {
		return model.ResolveResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error creating registry entry: %w", err))
	}

	if u.amqpQueue != nil {
		if err := u.sendNewRegistryMessage(registryId.String(), identityProvider); err != nil {
			slog.Warn("error publishing new registry entry", slog.String("error", err.Error()))
		}
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

func (u *User) resolveAccountWithIdentityProvider(email string, provider *string) (model.ResolveResponse, error) {
	userRecord, err := u.userDb.GetUserByEmail(email)
	if err != nil {
		return model.ResolveResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting user by email: %w", err))
	}
	token, profile, err := u.loginValidUser(userRecord, provider)
	if err != nil {
		return model.ResolveResponse{}, err
	}
	return model.ResolveResponse{
		NextAuthStep: constants.SessionStep,
		Metadata: model.ResolveWithProviderMetadata{
			Token:   token,
			Profile: profile,
		},
	}, nil
}

func (u *User) ResolveUserEmail(email string, identityProvider *string) (model.ResolveResponse, error) {
	if valErrs, err := u.userValidator.ValidateEmail(email); err != nil {
		return model.ResolveResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error validating mail: %w", err))
	} else if len(valErrs) > 0 {
		return model.ResolveResponse{}, app_errors.NewAppValidationError(valErrs)
	}

	hasAccount, err := u.checkIfEmailHasAccount(email)
	if err != nil {
		return model.ResolveResponse{}, err
	}

	if hasAccount {
		slog.Info("user email resolved successfully: it has account", slog.String("email", email))
		if identityProvider != nil {
			return u.resolveAccountWithIdentityProvider(email, identityProvider)
		}
		return model.ResolveResponse{
			NextAuthStep: constants.LoginStep,
			Metadata:     nil,
		}, nil
	}

	exists, err := u.registryDb.CheckIfRegistryEntryExistsByEmail(email)
	if err != nil {
		return model.ResolveResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if registry entry exists: %w", err))
	}

	if exists {
		slog.Info("user email resolved successfully: it has registry entry", slog.String("email", email))
		return u.resolveExistingRegistry(email)
	}

	slog.Info("user email resolved successfully: it doesnt have account", slog.String("email", email))
	return u.createNewRegistry(email, identityProvider)
}
