package service

import (
	"fmt"
	"net/http"
	"log/slog"
	"github.com/google/uuid"
	"users-service/src/model"
	"users-service/src/constants"
	"users-service/src/app_errors"
	"users-service/src/database/register_options"
)

func (u *User) GetRegisterOptions() map[string]interface{} {
	slog.Info("register options retrieved successfully")

	locations := []model.Location{}
	for id, name := range register_options.GetAllLocationsAndIds() {
		locations = append(locations, model.Location{Id: id, Name: name})
	}

	interests := []model.Interest{}
	for id, interest := range register_options.GetAllInterestsAndIds() {
		interests = append(interests, model.Interest{Id: id, Interest: interest})
	}

	return map[string]interface{}{
		"locations": locations,
		"interests": interests,
	}
}

func (u *User) validateRegistryEntry(id uuid.UUID) error {
	slog.Info("checking if registry entry exists")

	hasRegistry, err := u.registryDb.CheckIfRegistryEntryExists(id)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if registry entry exists: %w", err))
	}

	if !hasRegistry {
		return app_errors.NewAppError(http.StatusNotFound, RegistryNotFound, ErrRegistryNotFound)
	}

	return nil
}

func (u *User) validateRegistryStep(id uuid.UUID, step string) error {
	registry, err := u.registryDb.GetRegistryEntry(id)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting registry entry: %w", err))
	}

	actual_step := getStepForRegistryEntry(registry)
	if actual_step != step {
		return app_errors.NewAppError(http.StatusConflict, InvalidRegistryStep, fmt.Errorf("invalid registry step, should be %s, it is %s", actual_step, step))
	}

	return nil
}


func (u *User) SendVerificationEmail(id uuid.UUID) error {
	slog.Info("sending verification email")

	if err := u.validateRegistryEntry(id); err != nil {
		return err
	}

	if err := u.validateRegistryStep(id, constants.EmailVerificationStep); err != nil {
		return err
	}

	// if err := u.registryDb.SendVerificationEmail(id, email); err != nil {
	// 	return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error sending verification email: %w", err))
	// }

	slog.Info("verification email sent successfully", slog.String("registration_id", id.String()))
	return nil
}

func (u *User) VerifyEmail(id uuid.UUID, pin string) error {
	slog.Info("verifying email")

	if err := u.validateRegistryEntry(id); err != nil {
		return err
	}

	if err := u.validateRegistryStep(id, constants.EmailVerificationStep); err != nil {
		return err
	}

	if err := u.registryDb.VerifyEmail(id); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error verifying email: %w", err))
	}

	slog.Info("email verified successfully", slog.String("registration_id", id.String()))
	return nil
}

func (u *User) AddPersonalInfo(id uuid.UUID, data model.UserPersonalInfoRequest) error {
	slog.Info("adding personal info")

	if valErrs, err := u.userValidator.ValidatePersonalInfo(data); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error validating user personal info: %w", err))
	} else if len(valErrs) > 0 {
		return app_errors.NewAppValidationError(valErrs)
	}

	if err := u.validateRegistryEntry(id); err != nil {
		return err
	}

	if err := u.validateRegistryStep(id, constants.PersonalInfoStep); err != nil {
		return err
	}

	userInfo, err := generateUserPersonalInfoRecordFromRequest(data)
	if err != nil {
		return err
	}

	if err := u.registryDb.AddPersonalInfoToRegistryEntry(id, *userInfo); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error adding personal info to registry entry: %w", err))
	}

	slog.Info("personal info added successfully", slog.String("registration_id", id.String()))
	return nil
}

func (u *User) AddInterests(id uuid.UUID, interestsIds []int) error {
	slog.Info("adding interests")

	if valErrs, err := u.userValidator.ValidateInterests(interestsIds); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error validating user personal info: %w", err))
	} else if len(valErrs) > 0 {
		return app_errors.NewAppValidationError(valErrs)
	}

	if err := u.validateRegistryEntry(id); err != nil {
		return err
	}

	if err := u.validateRegistryStep(id, constants.InterestsStep); err != nil {
		return err
	}

	interestsNames, err := extractInterestNames(interestsIds)
	if err != nil {
		return app_errors.NewAppError(http.StatusBadRequest, InvalidInterest, fmt.Errorf("error extracting interest names: %w", err))
	}

	if err := u.registryDb.AddInterestsToRegistryEntry(id, interestsNames); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error adding interests to registry entry: %w", err))
	}

	slog.Info("interests added successfully", slog.String("registration id", id.String()))
	return nil
}

func (u *User) createUserWithInterestsFromRegistry(registry model.RegistryEntry, interestsNames []string) (model.UserResponse, error) {
	userRecord := generateUserRecordFromRegistryEntry(registry)
	createdUser, err := u.userDb.CreateUser(userRecord)
	if err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error creating user: %w", err))
	}

	err = u.interestDb.AssociateInterestsToUser(createdUser.Id, interestsNames)
	if err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error associating interest to user: %w", err))
	}

	return createUserResponseFromUserRecordAndInterests(createdUser, interestsNames), nil
}


func (u *User) CompleteRegistry(id uuid.UUID) (model.UserResponse, error) {
	slog.Info("completing registry")

	if err := u.validateRegistryStep(id, constants.CompleteStep); err != nil {
		return model.UserResponse{}, err
	}

	registry, err := u.registryDb.GetRegistryEntry(id)
	if err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting registry entry: %w", err))
	}

	if err := u.registryDb.DeleteRegistryEntry(id); err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error deleting registry entry: %w", err))
	}

	userResponse, err := u.createUserWithInterestsFromRegistry(registry, registry.Interests)
	if err != nil {
		return model.UserResponse{}, err
	}

	slog.Info("registry completed successfully", slog.String("registration id", id.String()))
	return userResponse, nil
}