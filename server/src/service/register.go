package service

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"users-service/src/app_errors"
	"users-service/src/constants"
	"users-service/src/database"
	"users-service/src/database/register_options"
	"users-service/src/model"

	"github.com/google/uuid"
)

func (u *User) GetLocations() map[string]interface{} {
	locations := []model.Location{}
	for id, name := range register_options.GetAllLocationsAndIds() {
		locations = append(locations, model.Location{Id: id, Name: name})
	}

	slog.Info("locations retrieved successfully")
	return map[string]interface{}{
		"locations": locations,
	}
}
func (u *User) GetInterests() map[string]interface{} {

	interests := []model.Interest{}
	for id, interest := range register_options.GetAllInterestsAndIds() {
		interests = append(interests, model.Interest{Id: id, Interest: interest})
	}

	slog.Info("interests retrieved successfully")
	return map[string]interface{}{
		"interests": interests,
	}
}

func (u *User) validateRegistryEntryExists(id uuid.UUID) error {
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

// GenerateRandomInRange generates a random number in the range [low, hi)
func GenerateRandomInRange(low, hi int) int {
    return low + rand.Intn(hi-low)
}

func (u *User) SendVerificationEmail(id uuid.UUID) error {
	slog.Info("sending verification email")

	registry, err := u.registryDb.GetRegistryEntry(id)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return app_errors.NewAppError(http.StatusNotFound, RegistryNotFound, ErrRegistryNotFound)
		}
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting registry entry: %w", err))
	}

	actual_step := getStepForRegistryEntry(registry)
	if actual_step != constants.EmailVerificationStep {
		return app_errors.NewAppError(http.StatusConflict, InvalidRegistryStep, fmt.Errorf("invalid registry step, should be %s, it is %s", actual_step, constants.EmailVerificationStep))
	}

	code := strconv.Itoa(GenerateRandomInRange(100000, 999999)) // random 6 digit number

	if err := u.registryDb.SetEmailVerificationPin(id, code); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error setting email verification pin: %w", err))
	}

	if err := SendVerificationEmail(registry.Email, code); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error sending verification email: %w", err))
	}

	slog.Info("verification email sent successfully", slog.String("registration_id", id.String()))
	return nil
}

func (u *User) VerifyEmail(id uuid.UUID, pin string) error {
	slog.Info("verifying email")

	if err := u.validateRegistryEntryExists(id); err != nil {
		return err
	}

	if err := u.validateRegistryStep(id, constants.EmailVerificationStep); err != nil {
		return err
	}

	verificationPin, err := u.registryDb.GetEmailVerificationPin(id)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return app_errors.NewAppError(http.StatusNotFound, VerificationPinNotFound, ErrVerificationPinNotFound)
		}
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting email verification pin: %w", err))
	}

	if verificationPin != pin {
		return app_errors.NewAppError(http.StatusBadRequest, VerificationPinNotFound, ErrVerificationPinNotFound)
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

	if err := u.validateRegistryEntryExists(id); err != nil {
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

	if err := u.validateRegistryEntryExists(id); err != nil {
		return err
	}

	if err := u.validateRegistryStep(id, constants.InterestsStep); err != nil {
		return err
	}

	interestsNames := extractInterestNamesFromValidIds(interestsIds)
	if err := u.registryDb.AddInterestsToRegistryEntry(id, interestsNames); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error adding interests to registry entry: %w", err))
	}

	slog.Info("interests added successfully", slog.String("registration id", id.String()))
	return nil
}

func (u *User) createUserFromRegistry(registry model.RegistryEntry) (model.UserPrivateProfile, error) {
	userRecord := generateUserRecordFromRegistryEntry(registry)
	createdUser, err := u.userDb.CreateUser(userRecord)
	if err != nil {
		return model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error creating user: %w", err))
	}

	return u.createUserPrivateProfileFromUserRecord(createdUser)
}

func (u *User) CompleteRegistry(id uuid.UUID) (model.UserPrivateProfile, error) {
	slog.Info("completing registry")

	if err := u.validateRegistryStep(id, constants.CompleteStep); err != nil {
		return model.UserPrivateProfile{}, err
	}

	registry, err := u.registryDb.GetRegistryEntry(id)
	if err != nil {
		return model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting registry entry: %w", err))
	}

	if err := u.registryDb.DeleteRegistryEntry(id); err != nil {
		return model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error deleting registry entry: %w", err))
	}

	userResponse, err := u.createUserFromRegistry(registry)
	if err != nil {
		return model.UserPrivateProfile{}, err
	}

	slog.Info("registry completed successfully", slog.String("registration id", id.String()))
	return userResponse, nil
}
