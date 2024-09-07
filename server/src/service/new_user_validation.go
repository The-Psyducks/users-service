package service

import (
	"regexp"
	"users-service/src/database"
	"users-service/src/model"

	"github.com/go-playground/validator/v10"
)

type UserCreationValidator struct {
	validationErrors []model.ValidationError
}

func NewUserCreationValidator() *UserCreationValidator {
	return &UserCreationValidator{}
}

func (u *UserCreationValidator) Validate(user model.UserRequest) []model.ValidationError {
	u.clearValidationErrors()
	validate := validator.New()

	customValidators := map[string]validator.Func{
		"usernamevalidator":  u.usernameValidator,
		"passwordvalidator":  u.passwordValidator,
		"mailvalidator":      u.mailValidator,
		"locationvalidator":  u.locationValidator,
		"interestsvalidator": u.interestsValidator,
	}

	for name, validatorFunc := range customValidators {
		validate.RegisterValidation(name, validatorFunc)
	}

	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if _, isCustom := customValidators[err.ActualTag()]; !isCustom {
				u.addValidationError(err.Field(), err.Tag())
			}
		}
	}

	return u.validationErrors
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

func (u *UserCreationValidator) mailValidator(fl validator.FieldLevel) bool {
	mail := fl.Field().String()
	mailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(mailPattern, mail)
	if err != nil || !matched {
		u.addValidationError("Mail", "Invalid email format")
		return false
	}
	return true
}

func (u *UserCreationValidator) usernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 4 || len(username) > 20 {
		u.addValidationError("Username", "Username must be between 4 and 20 characters long")
		return false
	}
	return true
}

func (u *UserCreationValidator) passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 || len(password) > 20 {
		u.addValidationError("Password", "Password must be between 8 and 20 characters long")
		return false
	}

	patterns := []string{`[A-Z]`, `[a-z]`, `[0-9]`, `[!@#$%^&*()]`}
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, password)
		if err != nil || !matched {
			u.addValidationError("Password", "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
			return false
		}
	}
	return true
}

func (u *UserCreationValidator) locationValidator(fl validator.FieldLevel) bool {
	location := fl.Field().Int()
	if database.GetLocationName(int32(location)) == "" {
		u.addValidationError("Location", "Invalid location")
		return false
	}
	return true
}

func (u *UserCreationValidator) interestsValidator(fl validator.FieldLevel) bool {
	interests := fl.Field().Interface().([]int32)
	for _, interest := range interests {
		if database.GetInterestName(interest) == "" {
			u.addValidationError("Interests", "Invalid interest")
			return false
		}
	}
	return true
}
