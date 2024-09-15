package service

import (
	"fmt"
	"log/slog"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/constants"
	"users-service/src/database/register_options"
	"users-service/src/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func createUserResponseFromUserRecordAndInterests(record model.UserRecord, interests []string) model.UserResponse {
	return model.UserResponse{
		Id:        record.Id,
		UserName:  record.UserName,
		FirstName: record.FirstName,
		LastName:  record.LastName,
		Mail:      record.Mail,
		Location:  record.Location,
		Interests: interests,
	}
}

func generateUserPersonalInfoRecordFromRequest(request model.UserPersonalInfoRequest) (*model.UserPersonalInfoRecord, error) {
	password, err := hashPassword(request.Password)
	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, "Internal server error", fmt.Errorf("error hashing password: %w", err))
	}

	return &model.UserPersonalInfoRecord{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		UserName:  request.UserName,
		Password:  password,
		Location:  register_options.GetLocationName(request.LocationId),
	}, nil
}

func generateUserRecordFromRegistryEntry(registry model.RegistryEntry) model.UserRecord {
	return model.UserRecord{
		UserName:  registry.PersonalInfo.UserName,
		FirstName: registry.PersonalInfo.FirstName,
		LastName:  registry.PersonalInfo.LastName,
		Mail:      registry.Email,
		Password:  registry.PersonalInfo.Password,
		Location:  registry.PersonalInfo.Location,
	}
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func extractInterestNames(interests []int) ([]string, error) {
	interestsNames := make([]string, len(interests))
	for i, interest := range interests {
		if name := register_options.GetInterestName(interest); name != "" {
			interestsNames[i] = name
		} else {
			return nil, fmt.Errorf("error: interest with id %d not found", interest)
		}

	}
	return interestsNames, nil
}

func (u *User) checkIfEmailHasAccount(email string) (bool, error) {
	slog.Info("checking if email has account")

	user, err := u.userDb.CheckIfMailExists(email)

	if err != nil {
		return false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if email exists: %w", err))
	}

	return user, nil

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

func getStepForRegistryEntry(entry model.RegistryEntry) string {
	if !entry.EmailVerified {
		return constants.EmailVerificationStep
	}

	if entry.PersonalInfo.FirstName == "" {
		return constants.PersonalInfoStep
	}

	if len(entry.Interests) == 0 {
		return constants.InterestsStep
	}

	return constants.CompleteStep
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
