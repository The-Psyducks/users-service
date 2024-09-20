package controller

import (
	"fmt"
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

func (u *User) HandleNoRoute(c *gin.Context) {
	err := app_errors.NewAppError(http.StatusMethodNotAllowed, "Method not found", fmt.Errorf("route not found"))
	_ = c.Error(err)
}

func (u *User) GetLocations(c *gin.Context) {
	data := u.service.GetLocations()

	c.JSON(http.StatusOK, data)
}

func (u *User) GetInterests(c *gin.Context) {
	data := u.service.GetInterests()

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

	c.JSON(http.StatusNoContent, gin.H{})
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

func (u *User) Login(c *gin.Context) {
	var data model.UserLoginRequest

	if err := c.BindJSON(&data); err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	token, err := u.service.CheckLoginCredentials(data)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

func (u *User) GetUserProfile(c *gin.Context) {
	username := c.Param("username")
	userSessionId := c.GetString("session_user_id")
	if userSessionId == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		_ = c.Error(err)
		return
	}

	user, err := u.service.GetUserProfile(userSessionId, username)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *User) FollowUser(c *gin.Context) {
	userSessionId := c.GetString("session_user_id")
	if userSessionId == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		_ = c.Error(err)
		return
	}

	var userToFollow struct {
		Username string `json:"username" validate:"required"`
	}
	if err := c.BindJSON(&userToFollow); err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	err := u.service.FollowUser(userSessionId, userToFollow.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (u *User) UnfollowUser(c *gin.Context) {
	userSessionId := c.GetString("session_user_id")
	if userSessionId == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		_ = c.Error(err)
		return
	}

	var userToUnfollow struct {
		Username string `json:"username" validate:"required"`
	}
	if err := c.BindJSON(&userToUnfollow); err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	err := u.service.UnfollowUser(userSessionId, userToUnfollow.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (u *User) GetFollowers(c *gin.Context) {
	username := c.Param("username")
	userSessionId := c.GetString("session_user_id")
	if userSessionId == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		_ = c.Error(err)
		return
	}
	followers, err := u.service.GetFollowers(username, userSessionId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, model.FollowersResponse{Followers: followers})
}

func (u *User) GetFollowing(c *gin.Context) {
	username := c.Param("username")
	userSessionId := c.GetString("session_user_id")
	if userSessionId == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		_ = c.Error(err)
		return
	}

	following, err := u.service.GetFollowing(username, userSessionId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, model.FollowingResponse{Following: following})
}