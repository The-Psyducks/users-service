package service

import (
	"errors"
	"net/http"
	"log/slog"
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

func CreateUserResponseFromUserRecordAndInterests(record model.UserRecord, interests []string) model.UserResponse {
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

func CreateUserRecordFromUserRequest(req *model.UserRequest) *model.UserRecord {
	return &model.UserRecord{
		UserName:  req.UserName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Mail:      req.Mail,
		Password:  req.Password,
		Location:  database.GetLocationName(req.Location),
	}
}

func (u *User) CreateUser(data model.UserRequest) (model.UserResponse, error) {
	//validate data
	userRecord := CreateUserRecordFromUserRequest(&data)
	createdUser, err := u.user_db.CreateUser(*userRecord)

	if err != nil {
		//todo: user already exists
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, "error creating user", err)
	}
	
	interests, err := u.interest_db.AssociateInterestsToUser(createdUser.Id, data.Interests)
	
	if err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, "error associating interest to user", err)
	}

	interestsNames := make([]string, len(interests))
	for i, interest := range interests {
		interestsNames[i] = interest.Name
	}
	slog.Info("user created succesfully", slog.String("user_id", createdUser.Id.String()))
	return CreateUserResponseFromUserRecordAndInterests(createdUser, interestsNames), nil
}

func (u *User) GetRegisterOptions() map[string]interface{} {
	slog.Info("register optiones retrieved succesfully")
	return map[string]interface{}{
		"locations": database.GetAllLocations(),
		"interests": database.GetAllInterests(),
	}
}

func (u *User) GetUserById(id string) (model.UserResponse, error) {
	userRecord, err := u.user_db.GetUserById(id)

	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return model.UserResponse{}, app_errors.NewAppError(http.StatusNotFound, "user not found", err)
		}
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, "error retrieving user", err)
	}

	interests, err := u.interest_db.GetInterestsNamesForUserId(userRecord.Id)

	if err != nil {
		return model.UserResponse{}, app_errors.NewAppError(http.StatusInternalServerError, "error getting interests from user", err)
	}

	slog.Info("user retrieved succesfully", slog.String("user_id", userRecord.Id.String()))
	return CreateUserResponseFromUserRecordAndInterests(userRecord, interests), nil
}
