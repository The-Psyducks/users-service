package app_errors

import (
	"fmt"
	"net/http"
	"strings"
	"users-service/src/model"
)

type AppValidationError struct {
	Code    int
	Message string
	Errors  []model.ValidationError
}

func (e *AppValidationError) Error() string {
	var errorMessages []string
	for _, err := range e.Errors {
		errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return fmt.Sprintf("Validation errors: %s", strings.Join(errorMessages, "; "))
}

func NewAppValidationError(errors []model.ValidationError) *AppValidationError {
	return &AppValidationError{
		Code:    http.StatusBadRequest,
		Message: "Validation error",
		Errors:  errors,
	}
}
