package controller

import (
	"fmt"
	"net/http"
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"users-service/src/app_errors"
	"users-service/src/constants"
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

func (u *User) GetUserProfileByUsername(c *gin.Context) {
	username := c.Param("username")
	userSessionId := c.GetString("session_user_id")
	if userSessionId == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		_ = c.Error(err)
		return
	}

	user, err := u.service.GetUserProfileByUsername(userSessionId, username)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *User) GetUserProfileById(c *gin.Context) {
	idString := c.Param("id")
	userSessionId := c.GetString("session_user_id")
	if userSessionId == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		_ = c.Error(err)
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	user, err := u.service.GetUserProfileById(userSessionId, id)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *User) FollowUser(c *gin.Context) {
	username := c.Param("username")
	userSessionId := c.GetString("session_user_id")
	if userSessionId == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		_ = c.Error(err)
		return
	}

	err := u.service.FollowUser(userSessionId, username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (u *User) UnfollowUser(c *gin.Context) {
	username := c.Param("username")
	userSessionId := c.GetString("session_user_id")
	if userSessionId == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		_ = c.Error(err)
		return
	}

	err := u.service.UnfollowUser(userSessionId, username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func getPaginationParams(c *gin.Context) (string, int, int, error) {
	timestampStr := c.DefaultQuery("time", time.Now().UTC().Format(time.RFC3339))
	_, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid 'timestamp' value in request. Must be in RFC3339 format.", err)
		return "", 0, 0, err
	}

	skipStr := c.DefaultQuery("skip", "0")
	skipInt, err := strconv.Atoi(skipStr)
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid 'skip' value in request", err)
		return "", 0, 0, err
	}

	limitStr := c.DefaultQuery("limit", "20")
	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid 'limit' value in request", err)
		return "", 0, 0, err
	}

	if limitInt > constants.MaxPaginationLimit {
		slog.Warn("limit is higher than the max pagination limit, using default", 
			slog.Int("limit", limitInt), slog.Int("max", constants.MaxPaginationLimit))
		limitInt = constants.MaxPaginationLimit
	}
	return timestampStr, skipInt, limitInt, nil
}

func (u *User) GetFollowers(c *gin.Context) {
	timestamp, skip, limit, err := getPaginationParams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	username := c.Param("username")
	userSessionId := c.GetString("session_user_id")
	if userSessionId == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		_ = c.Error(err)
		return
	}
	followers, hasMore, err := u.service.GetFollowers(username, userSessionId, timestamp, skip, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := model.FollowersPaginationResponse{
		Followers: followers,
		Pagination: model.Pagination{
			Limit: limit,
		},
	}
	if hasMore {
		response.Pagination.NextOffset = skip + limit
	}
	c.JSON(http.StatusOK, response)
}

func (u *User) GetFollowing(c *gin.Context) {
	timestamp, skip, limit, err := getPaginationParams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	username := c.Param("username")
	userSessionId := c.GetString("session_user_id")
	if userSessionId == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		_ = c.Error(err)
		return
	}

	following, hasMore, err := u.service.GetFollowing(username, userSessionId, timestamp, skip, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := model.FollowingPaginationResponse{
		Following: following,
		Pagination: model.Pagination{
			Limit: limit,
		},
	}
	if hasMore {
		response.Pagination.NextOffset = skip + limit
	}
	c.JSON(http.StatusOK, response)
}
