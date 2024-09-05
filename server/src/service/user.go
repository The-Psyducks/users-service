package service

import (
	"users-service/src/database"
	"users-service/src/model"
)

type User struct {
	db database.Database
}

func UserRecordToUserResponse(record model.UserRecord) model.UserResponse {
	return model.UserResponse{
		Id: record.Id,
		UserName: record.UserName,
		Name: record.Name,
		Mail: record.Mail,
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

func NewUser(db database.Database) *User {
	return &User{db: db}
}