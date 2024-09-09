package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

func (u *User) CreateUser(c *gin.Context) {
	var data model.UserRequest

	if err := c.BindJSON(&data); err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	user, err := u.service.CreateUser(data)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (u *User) GetRegisterOptions(c *gin.Context) {
	data := u.service.GetRegisterOptions()

	c.JSON(http.StatusOK, data)
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
