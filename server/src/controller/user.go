package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"users-service/src/app_errors"
	"users-service/src/model"
	"users-service/src/service"
)

type User struct {
	service *service.User
}

func CreateUserController(service *service.User) *User {
	return &User{service: service}
}

func (u *User) GetRegisterOptions(c *gin.Context) {
	data := u.service.GetRegisterOptions()

	c.JSON(http.StatusOK, data)
}

func (u *User) ResolveUserEmail(c *gin.Context) {
	var data model.ResolveRequest

	if err := c.BindJSON(&data); err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	user, err := u.service.ResolveUserEmail(data)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *User) SendVerificationEmail(c *gin.Context) {
	registrationId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	err = u.service.SendVerificationEmail(registrationId)

	if err != nil {
	_ = c.Error(err)
	return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (u *User) VerifyEmail(c *gin.Context) {
	var verificationRequest struct {
		Pin string `json:"pin" validate:"required"`
	}

	if err := c.BindJSON(&verificationRequest); err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	registrationId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	err = u.service.VerifyEmail(registrationId, verificationRequest.Pin)
	if err != nil {
	_ = c.Error(err)
	return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (u *User) AddPersonalInfo(c *gin.Context) {
	var data model.UserPersonalInfoRequest
	
	if err := c.BindJSON(&data); err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}
	
	registrationId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	err = u.service.AddPersonalInfo(registrationId, data)

	if err != nil {
	_ = c.Error(err)
	return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (u *User) AddInterests(c *gin.Context) {
	var interests struct {
		InterestsIds []int `json:"interests" validate:"required"`
	}
	
	if err := c.BindJSON(&interests); err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}
	
	registrationId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	err = u.service.AddInterests(registrationId, interests.InterestsIds)

	if err != nil {
	_ = c.Error(err)
	return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (u *User) CompleteRegistry(c *gin.Context) {
	registrationId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}
	
	userResponse, err := u.service.CompleteRegistry(registrationId)
	if err != nil {
	_ = c.Error(err)
	return
	}

	c.JSON(http.StatusOK, userResponse)
}

func (u *User) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := u.service.GetUserByUsername(username)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *User) Login(c *gin.Context) {
	var data model.UserLoginRequest

	if err := c.BindJSON(&data); err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	valid, err := u.service.CheckLoginCredentials(data)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": valid})
}
