package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"users-service/src/service"
	"users-service/src/model"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.service.CreateUser(data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (u *User) GetUserById(c *gin.Context) {
	id := c.Param("id")

	user, err := u.service.GetUserById(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}