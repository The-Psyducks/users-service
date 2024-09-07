package service

import (
	"regexp"
	"users-service/src/database"
	"users-service/src/model"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	FieldName string `json:"field_name"`
	Message   string `json:"message"`
}

type UserCreationValidator struct {
	user_db          database.UserDatabase
	validationErrors []ValidationError
}

func NewUserCreationValidator(user_db database.UserDatabase) *UserCreationValidator {
	return &UserCreationValidator{
		user_db: user_db,
	}
}

func (u *UserCreationValidator) Validate(user model.UserRequest) []ValidationError {
	u.clearValidationErrors()
	validate := validator.New()

    customValidators := map[string]validator.Func{
        "usernamevalidator":   u.usernameValidator,
        "passwordvalidator":   u.passwordValidator,
        "mailvalidator":       u.mailValidator,
        "locationvalidator":   u.locationValidator,
        "interestsvalidator":  u.interestsValidator,
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
	u.validationErrors = []ValidationError{}
}

func (u *UserCreationValidator) addValidationError(fieldName, message string) {
	u.validationErrors = append(u.validationErrors, ValidationError{
		FieldName: fieldName,
		Message:   message,
	})
}

func (u *UserCreationValidator) usernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 4 || len(username) > 20 {
		u.addValidationError("Username", "Username must be between 4 and 20 characters long")
		return false
	}

	exists, err := u.user_db.CheckIfUsernameExists(username)
	if err != nil {
		u.addValidationError("UserName", "Error checking username availability")
		return false
	}
	if exists {
		u.addValidationError("UserName", "This username is already taken")
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

func (u *UserCreationValidator) mailValidator(fl validator.FieldLevel) bool {
	mail := fl.Field().String()
	mailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(mailPattern, mail)
	if err != nil || !matched {
		u.addValidationError("Mail", "Invalid email format")
		return false
	}

	exists, err := u.user_db.CheckIfMailExists(mail)
	if err != nil {
		u.addValidationError("Mail", "Error checking email availability")
		return false
	}
	if exists {
		u.addValidationError("Mail", "This email is already taken")
		return false
	}
	return true
}
