package controller

import (
	"fmt"
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
		fmt.Println(err)
		err = app_errors.NewAppError(http.StatusBadRequest, "invalid request", err)
		c.Error(err)
		return
	}

	user, err := u.service.CreateUser(data)

	if err != nil {
		c.Error(err)
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
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}
