package service

import (
	"fmt"
	"log/slog"
	"regexp"
	"users-service/src/constants"
	"users-service/src/database/register_options"
	"users-service/src/database/users_db"
	"users-service/src/model"

	"github.com/go-playground/validator/v10"
)

type UserCreationValidator struct {
	usersDb          users_db.UserDatabase
	validationErrors []model.ValidationError
}

func NewUserCreationValidator(usersDb users_db.UserDatabase) *UserCreationValidator {
	return &UserCreationValidator{
		usersDb: usersDb,
	}
}

func (u *UserCreationValidator) ValidateUpdatedPrivateProfile(updatedProfile model.UpdateUserPrivateProfileRequest) ([]model.ValidationError, error) {
	u.clearValidationErrors()

	validate := validator.New()

	customValidators := map[string]validator.Func{
		"firstnamevalidator": u.firstnamevalidator,
		"lastnamevalidator":  u.lastnamevalidator,
		"usernamevalidator":  u.usernameValidator,
		"locationvalidator":  u.locationValidator,
		"interestsvalidator": u.interestsValidator,
	}

	for name, validatorFunc := range customValidators {
		if err := validate.RegisterValidation(name, validatorFunc); err != nil {
			slog.Error("Error registering custom validator", slog.String("error: ", err.Error()))
			return []model.ValidationError{}, err
		}
	}

	err := validate.Struct(updatedProfile)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if _, isCustom := customValidators[err.ActualTag()]; !isCustom {
				u.addValidationError(err.Field(), err.Tag())
			}
		}
	}

	return u.validationErrors, nil
}

func (u *UserCreationValidator) ValidateEmail(email string) ([]model.ValidationError, error) {
	u.clearValidationErrors()

	validate := validator.New()
	if err := validate.RegisterValidation("emailvalidator", u.emailValidator); err != nil {
		slog.Error("Error registering custom validator", slog.String("error: ", err.Error()))
		return []model.ValidationError{}, err
	}

	err := validate.Var(email, "emailvalidator")
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.ActualTag() != "emailvalidator" {
				u.addValidationError("interests", err.Error())
			}
		}
	}

	return u.validationErrors, nil
}

func (u *UserCreationValidator) ValidatePersonalInfo(personalInfo model.UserPersonalInfoRequest) ([]model.ValidationError, error) {
	u.clearValidationErrors()
	validate := validator.New()

	customValidators := map[string]validator.Func{
		"firstnamevalidator": u.firstnamevalidator,
		"lastnamevalidator":  u.lastnamevalidator,
		"usernamevalidator":  u.usernameValidator,
		"passwordvalidator":  u.passwordValidator,
		"locationvalidator":  u.locationValidator,
	}

	for name, validatorFunc := range customValidators {
		if err := validate.RegisterValidation(name, validatorFunc); err != nil {
			slog.Error("Error registering custom validator", slog.String("error: ", err.Error()))
			return []model.ValidationError{}, err
		}
	}

	err := validate.Struct(personalInfo)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if _, isCustom := customValidators[err.ActualTag()]; !isCustom {
				u.addValidationError(err.Field(), err.Tag())
			}
		}
	}

	return u.validationErrors, nil
}

func (u *UserCreationValidator) ValidateInterests(interests []int) ([]model.ValidationError, error) {
	u.clearValidationErrors()
	validate := validator.New()

	if err := validate.RegisterValidation("interestsvalidator", u.interestsValidator); err != nil {
		slog.Error("Error registering custom validator", slog.String("error: ", err.Error()))
		return []model.ValidationError{}, err
	}

	err := validate.Var(interests, "interestsvalidator")
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.ActualTag() != "interestsvalidator" {
				u.addValidationError("interests", err.Error())
			}
		}
	}

	return u.validationErrors, nil
}

func (u *UserCreationValidator) clearValidationErrors() {
	u.validationErrors = []model.ValidationError{}
}

func (u *UserCreationValidator) addValidationError(fieldName, message string) {
	u.validationErrors = append(u.validationErrors, model.ValidationError{
		Field:   fieldName,
		Message: message,
	})
}

func (u *UserCreationValidator) emailValidator(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	if len(email) < constants.MinEmailLength || len(email) > constants.MaxEmailLength {
		u.addValidationError("email", fmt.Sprintf("Email must be between %d and %d characters long", constants.MinEmailLength, constants.MaxEmailLength))
		return false
	}

	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailPattern, email)
	if err != nil || !matched {
		u.addValidationError("email", "Invalid email format")
		return false
	}
	return true
}

func (u *UserCreationValidator) firstnamevalidator(fl validator.FieldLevel) bool {
	first_name := fl.Field().String()
	if len(first_name) > constants.MaxFirstNameLength {
		u.addValidationError("first_name", fmt.Sprintf("First name must be less than %d characters long", constants.MaxFirstNameLength))
		return false
	}
	return true
}

func (u *UserCreationValidator) lastnamevalidator(fl validator.FieldLevel) bool {
	last_name := fl.Field().String()
	if len(last_name) > constants.MaxLastNameLength {
		u.addValidationError("last_name", fmt.Sprintf("Last name must be less than %d characters long", constants.MaxLastNameLength))
		return false
	}
	return true
}
func (u *UserCreationValidator) usernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < constants.MinUsernameLength || len(username) > constants.MaxUsernameLength {
		u.addValidationError("username", fmt.Sprintf("Username must be between %d and %d characters long", constants.MinUsernameLength, constants.MaxUsernameLength))
		return false
	}

	user, err := u.usersDb.CheckIfUsernameExists(username)
	if err != nil {
		u.addValidationError("username", "Error checking if username exists")
		return false
	}

	if user {
		u.addValidationError("username", "Username already exists")
		return false
	}

	return true
}

func (u *UserCreationValidator) passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < constants.MinPasswordLength || len(password) > constants.MaxPasswordLength {
		u.addValidationError("password", "Password must be between 8 and 20 characters long")
		return false
	}

	patterns := []string{`[A-Z]`, `[a-z]`, `[0-9]`, `[!@#$%^&*()]`}
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, password)
		if err != nil || !matched {
			u.addValidationError("password", "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
			return false
		}
	}
	return true
}

func (u *UserCreationValidator) locationValidator(fl validator.FieldLevel) bool {
	location := fl.Field().Int()
	if register_options.GetLocationName(int(location)) == "" {
		u.addValidationError("location", "Invalid location")
		return false
	}
	return true
}

func (u *UserCreationValidator) interestsValidator(fl validator.FieldLevel) bool {
	interests := fl.Field().Interface().([]int)
	if len(interests) == 0 {
		u.addValidationError("interests", "A user must have at least one interest")
		return false
	}
	if len(interests) > constants.MaxInterests {
		u.addValidationError("interests", fmt.Sprintf("Interests must be less than %d", constants.MaxInterests))
		return false
	}

	seen := make(map[int]bool)
	for _, interest := range interests {
		if register_options.GetInterestName(interest) == "" {
			u.addValidationError("interests", "Invalid interest")
			return false
		}
		if seen[interest] {
			u.addValidationError("interests", "Duplicate interests")
			return false
		}
		seen[interest] = true
	}
	return true
}
