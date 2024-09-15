package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/constants"
	"users-service/src/database"
	"users-service/src/database/interests_db"
	"users-service/src/database/register_options"
	"users-service/src/database/registry_db"
	"users-service/src/database/users_db"
	"users-service/src/model"

	"github.com/google/uuid"
)

type User struct {
	userDb        users_db.UserDatabase
	interestDb    interests_db.InterestsDatabase
	registryDb    registry_db.RegistryDatabase
	userValidator *UserCreationValidator
}

func CreateUserService(userDb users_db.UserDatabase, interestDb interests_db.InterestsDatabase, registryDb registry_db.RegistryDatabase) *User {
	return &User{
		userDb:        userDb,
		interestDb:    interestDb,
		registryDb:    registryDb,
		userValidator: NewUserCreationValidator(userDb),
	}
}

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

func (u *User) ResolveUserEmail(data model.ResolveRequest) (model.ResolveResponse, error) {
	slog.Info("resolving user email")

	// chequeo de provider y verificacion del token

	if valErrs, err := u.userValidator.ValidateEmail(data.Email); err != nil {
		return model.ResolveResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error validating mail: %w", err))
	} else if len(valErrs) > 0 {
		return model.ResolveResponse{}, app_errors.NewAppValidationError(valErrs)
	}

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
	slog.Info("extracted interests names", slog.Any("ids", interestsIds), slog.Any("interests", interestsNames))
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error extracting interest names: %w", err))
	}

	if err := u.registryDb.AddInterestsToRegistryEntry(id, interestsNames); err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error adding interests to registry entry: %w", err))
	}

	slog.Info("interests added successfully", slog.String("registration id", id.String()))
	return nil
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

func (u *User) CheckLoginCredentials(data model.UserLoginRequest) (bool, error) {
	slog.Info("checking login information")

	userRecord, err := u.userDb.GetUserByUsername(data.UserName)

	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return false, app_errors.NewAppError(http.StatusNotFound, IncorrectUsernameOrPassword, errors.New("invalid username"))
		}
		return false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if !checkPasswordHash(data.Password, userRecord.Password) {
		return false, app_errors.NewAppError(http.StatusNotFound, IncorrectUsernameOrPassword, errors.New("invalid password"))
	}

	slog.Info("login information checked successfully", slog.String("username", userRecord.UserName))
	return true, nil
}

func (u *User) GetUserByUsername(username string) (model.UserResponse, error) {
	userRecord, err := u.userDb.GetUserByUsername(username)

	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return model.UserResponse{}, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	interests, err := u.interestDb.GetInterestsForUserId(userRecord.Id)
	if err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting interests from user: %w", err))
	}

	slog.Info("user retrieved succesfully", slog.String("username", userRecord.UserName))
	return createUserResponseFromUserRecordAndInterests(userRecord, interests), nil
}
