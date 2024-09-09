package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/database"
	"users-service/src/model"
)

type User struct {
	user_db     database.UserDatabase
	interest_db database.InterestsDatabase
}

func CreateUserService(user_db database.UserDatabase, interest_db database.InterestsDatabase) *User {
	return &User{
		user_db:     user_db,
		interest_db: interest_db,
	}
}

func (u *User) checkExistingUserData(username, mail string) *app_errors.AppError {
	usernameExists, err := u.user_db.CheckIfUsernameExists(username)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if username exists: %w", err))
	}
	if usernameExists {
		return app_errors.NewAppError(http.StatusConflict, UsernameOrMailAlreadyExists, errors.New("this username already exists"))
	}

	mailExists, err := u.user_db.CheckIfMailExists(mail)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if mail exists: %w", err))
	}
	if mailExists {
		return app_errors.NewAppError(http.StatusConflict, UsernameOrMailAlreadyExists, errors.New("this mail already exists"))
	}

	return nil
}

func (u *User) CreateUser(data model.UserRequest) (model.UserResponse, error) {
	slog.Info("creating new user")

	userValidator := NewUserCreationValidator()

	if valErrs, err := userValidator.Validate(data); err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error validating user: %w", err))
	} else if len(valErrs) > 0 {
		return model.UserResponse{}, app_errors.NewAppValidationError(valErrs)
	}

	if appErr := u.checkExistingUserData(data.UserName, data.Mail); appErr != nil {
		return model.UserResponse{}, appErr
	}

	userRecord, appErr := generateUserRecordFromUserRequest(&data)

	if appErr != nil {
		return model.UserResponse{}, nil
	}

	createdUser, err := u.user_db.CreateUser(*userRecord)

	if err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error creating user: %w", err))
	}

	interests, err := u.interest_db.AssociateInterestsToUser(createdUser.Id, data.InterestsIds)

	if err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error associating interest to user: %w", err))
	}

	interestsNames := extractInterestNames(interests)
	slog.Info("user created succesfully", slog.String("username", createdUser.UserName))
	return createUserResponseFromUserRecordAndInterests(createdUser, interestsNames), nil
}

func (u *User) GetRegisterOptions() map[string]interface{} {
	slog.Info("register options retrieved successfully")

	locations := []model.Location{}
	for id, name := range database.GetAllLocationsAndIds() {
		locations = append(locations, model.Location{Id: id, Name: name})
	}

	interests := []model.Interest{}
	for id, interest := range database.GetAllInterestsAndIds() {
		interests = append(interests, model.Interest{Id: id, Interest: interest})
	}

	return map[string]interface{}{
		"locations": locations,
		"interests": interests,
	}
}

func (u *User) CheckLoginCredentials(data model.UserLoginRequest) (bool, error) {
	slog.Info("checking login information")

	userRecord, err := u.user_db.GetUserByUsername(data.UserName)

	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return false, app_errors.NewAppError(http.StatusUnauthorized, IncorrectUsernameOrPassword, errors.New("invalid username"))
		}
		return false, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if !CheckPasswordHash(data.Password, userRecord.Password) {
		return false, app_errors.NewAppError(http.StatusUnauthorized, IncorrectUsernameOrPassword, errors.New("invalid password"))
	}

	slog.Info("login information checked successfully", slog.String("username", userRecord.UserName))
	return true, nil
}

func (u *User) GetUserByUsername(username string) (model.UserResponse, error) {
	userRecord, err := u.user_db.GetUserByUsername(username)

	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return model.UserResponse{}, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	interests, err := u.interest_db.GetInterestsNamesForUserId(userRecord.Id)

	if err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting interests from user: %w", err))
	}

	slog.Info("user retrieved succesfully", slog.String("username", userRecord.UserName))
	return createUserResponseFromUserRecordAndInterests(userRecord, interests), nil
}
