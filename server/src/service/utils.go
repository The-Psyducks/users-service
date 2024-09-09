package service

import (
	"fmt"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/database"
	"users-service/src/model"

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

func generateUserRecordFromUserRequest(req *model.UserRequest) (*model.UserRecord, *app_errors.AppError) {
	password, err := HashPassword(req.Password)

	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, "Internal server error", fmt.Errorf("error hashing password: %w", err))
	}

	return &model.UserRecord{
		UserName:  req.UserName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Mail:      req.Mail,
		Password:  password,
		Location:  database.GetLocationName(req.LocationId),
	}, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func extractInterestNames(interests []model.InterestRecord) []string {
	interestsNames := make([]string, len(interests))
	for i, interest := range interests {
		interestsNames[i] = interest.Name
	}
	return interestsNames
}