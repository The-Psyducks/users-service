package service

import (
	"fmt"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/constants"
	"users-service/src/database/register_options"
	"users-service/src/model"

	"golang.org/x/crypto/bcrypt"
)

func createUserResponseFromUserRecordAndInterests(record model.UserRecord, interests []string) model.UserResponse {
	return model.UserResponse{
		Id:        record.Id,
		UserName:  record.UserName,
		FirstName: record.FirstName,
		LastName:  record.LastName,
		Email:     record.Email,
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
		Email:     registry.Email,
		Password:  registry.PersonalInfo.Password,
		Location:  registry.PersonalInfo.Location,
	}
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", app_errors.NewAppError(http.StatusInternalServerError, "Internal server error", fmt.Errorf("error hashing password: %w", err))
	}
	
	return string(hashedPassword), nil
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
			return nil, fmt.Errorf("invalid interest: %d", interest)
		}

	}
	return interestsNames, nil
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

