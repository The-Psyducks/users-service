package service

import (
	"users-service/src/database"
	"users-service/src/model"
)

type User struct {
	db database.Database
}

func CreateUserService(db database.Database) *User {
	return &User{db: db}
}

func UserRecordToUserResponse(record model.UserRecord) model.UserResponse {
	return model.UserResponse{
		Id:       record.Id,
		UserName: record.UserName,
		Name:     record.Name,
		Mail:     record.Mail,
		Location: record.Location,
	}
}

func (u *User) CreateUser(data model.UserRequest) (model.UserResponse, error) {
	userRecord, err := u.db.CreateUser(data)

	//validate data

	if err != nil {
		return model.UserResponse{}, err
	}

	return UserRecordToUserResponse(userRecord), nil
}

func (u *User) GetUserById(id string) (model.UserResponse, error) {
	userRecord, err := u.db.GetUserById(id)

	if err != nil {
		return model.UserResponse{}, err
	}

	return UserRecordToUserResponse(userRecord), nil
}