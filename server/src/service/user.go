package service

import (
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

	interests := u.interest_db.AssociateInterestsToUser(createdUser.Id, data.Interests)
	interestsNames := make([]string, len(interests))
	for i, interest := range interests {
		interestsNames[i] = interest.Name
	}

	if err != nil {
		return model.UserResponse{}, err
	}

	return CreateUserResponseFromUserRecordAndInterests(createdUser, interestsNames), nil
}

func (u *User) GetRegisterOptions() map[string]interface{} {
	return map[string]interface{}{
		"locations": database.GetAllLocations(),
		"interests": database.GetAllInterests(),
	}
}

func (u *User) GetUserById(id string) (model.UserResponse, error) {
	userRecord, err := u.user_db.GetUserById(id)
	interests := u.interest_db.GetInterestsNamesForUserId(userRecord.Id)

	if err != nil {
		return model.UserResponse{}, err
	}

	return CreateUserResponseFromUserRecordAndInterests(userRecord, interests), nil
}
