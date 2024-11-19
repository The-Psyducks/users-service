package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"users-service/src/app_errors"
	"users-service/src/auth"
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

	var identityProvider *string
	switch data.ProviderData.Name {
	case constants.GoogleProvider:
		var googleMetadata model.GoogleAuthMetadata
		if err := json.Unmarshal(data.ProviderData.Metadata, &googleMetadata); err != nil {
			err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
			_ = c.Error(err)
			return
		}
		slog.Info("authenticating Google provider")
		if isValid, err := auth.IsGoogleTokenValid(googleMetadata.FirebaseTokenId); err != nil {
			err = app_errors.NewAppError(http.StatusInternalServerError, "Internal server error", fmt.Errorf("error validating google token: %w", err))
			_ = c.Error(err)
			return
		} else if !isValid {
			err := app_errors.NewAppError(http.StatusUnauthorized, "Invalid token", fmt.Errorf("google token is invalid: %s", googleMetadata.FirebaseTokenId))
			_ = c.Error(err)
			return
		}
		identityProvider = &data.ProviderData.Name
	case "":
	default:
		err := app_errors.NewAppError(http.StatusBadRequest, "Unknown provider type", nil)
		_ = c.Error(err)
		return
	}

	user, err := u.service.ResolveUserEmail(data.Email, identityProvider)
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
		Pin string `json:"pin" binding:"required"`
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

	token, profile, err := u.service.LoginUser(data)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token, "profile": profile})
}

func getSessionUserId(c *gin.Context) (uuid.UUID, error) {
	sessionUserIdString := c.GetString("session_user_id")
	if sessionUserIdString == "" {
		err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("session_user_id not found in context"))
		return uuid.Nil, err
	}
	sessionUserId, err := uuid.Parse(sessionUserIdString)
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid id in token", err)
		return uuid.Nil, err
	}
	return sessionUserId, nil
}

func getUrlIdAndSessionUserId(c *gin.Context) (uuid.UUID, uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		return uuid.Nil, uuid.Nil, err
	}

	sessionUserId, err := getSessionUserId(c)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	return id, sessionUserId, nil
}

func (u *User) GetUserProfileById(c *gin.Context) {
	id, userSessionId, err := getUrlIdAndSessionUserId(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userSessionIsAdmin := c.GetBool("session_user_admin")
	user, err := u.service.GetUserProfileById(userSessionId, userSessionIsAdmin, id)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *User) ModifyUserProfile(c *gin.Context) {
	sessionUserId, err := getSessionUserId(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var data model.UpdateUserPrivateProfileRequest
	if err := c.BindJSON(&data); err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	userProfile, err := u.service.ModifyUserProfile(sessionUserId, data)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, userProfile)
}

func (u *User) FollowUser(c *gin.Context) {
	userToFollowId, userSessionId, err := getUrlIdAndSessionUserId(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = u.service.FollowUser(userSessionId, userToFollowId, c.GetString("token"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (u *User) UnfollowUser(c *gin.Context) {
	userToUnfollowId, userSessionId, err := getUrlIdAndSessionUserId(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = u.service.UnfollowUser(userSessionId, userToUnfollowId)
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
	userToGetFollowersId, userSessionId, err := getUrlIdAndSessionUserId(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	timestamp, skip, limit, err := getPaginationParams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	followers, hasMore, err := u.service.GetFollowers(userToGetFollowersId, userSessionId, timestamp, skip, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := model.CreatePaginationResponse(followers, limit, skip, hasMore)
	c.JSON(http.StatusOK, response)
}

func (u *User) GetFollowing(c *gin.Context) {
	userToGetFollowingId, userSessionId, err := getUrlIdAndSessionUserId(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	timestamp, skip, limit, err := getPaginationParams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	following, hasMore, err := u.service.GetFollowing(userToGetFollowingId, userSessionId, timestamp, skip, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := model.CreatePaginationResponse(following, limit, skip, hasMore)
	c.JSON(http.StatusOK, response)
}

func (u *User) SearchUsers(c *gin.Context) {
	userSessionId, err := getSessionUserId(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	timestamp, skip, limit, err := getPaginationParams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	text := c.DefaultQuery("text", "")
	if strings.TrimSpace(text) == "" {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid 'text' value in request. Must not be empty.", fmt.Errorf("invalid search text"))
		_ = c.Error(err)
		return
	}

	users, hasMore, err := u.service.SearchUsers(userSessionId, text, timestamp, skip, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := model.CreatePaginationResponse(users, limit, skip, hasMore)
	c.JSON(http.StatusOK, response)
}

func (u *User) RecommendUsers(c *gin.Context) {
	userSessionId, err := getSessionUserId(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	timestamp, skip, limit, err := getPaginationParams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	users, hasMore, err := u.service.RecommendUsers(userSessionId, timestamp, skip, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := model.CreatePaginationResponse(users, limit, skip, hasMore)
	c.JSON(http.StatusOK, response)
}

func (u *User) GetAllUsers(c *gin.Context) {
	timestamp, skip, limit, err := getPaginationParams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userSessionIsAdmin := c.GetBool("session_user_admin")
	users, hasMore, err := u.service.GetAllUsers(userSessionIsAdmin, timestamp, skip, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := model.CreatePaginationResponse(users, limit, skip, hasMore)
	c.JSON(http.StatusOK, response)
}

func (u *User) GetUserInformation(c *gin.Context) {
	id, userSessionId, err := getUrlIdAndSessionUserId(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userSessionIsAdmin := c.GetBool("session_user_admin")
	user, err := u.service.GetUserInformation(userSessionId, userSessionIsAdmin, id)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *User) GetRegistrationMetrics(c *gin.Context) {
	userSessionIsAdmin := c.GetBool("session_user_admin")
	metrics, err := u.service.GetRegistrationMetrics(userSessionIsAdmin)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (u *User) GetLoginMetrics(c *gin.Context) {
	userSessionIsAdmin := c.GetBool("session_user_admin")
	metrics, err := u.service.GetLoginMetrics(userSessionIsAdmin)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (u *User) GetLocationMetrics(c *gin.Context) {
	userSessionIsAdmin := c.GetBool("session_user_admin")
	metrics, err := u.service.GetLocationMetrics(userSessionIsAdmin)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (u *User) GetUsersBlockedMetrics(c *gin.Context) {
	userSessionIsAdmin := c.GetBool("session_user_admin")
	metrics, err := u.service.GetUsersBlockedMetrics(userSessionIsAdmin)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (u *User) BlockUser(c *gin.Context) {
	userSessionIsAdmin := c.GetBool("session_user_admin")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}

	type blockRequest struct {
		Reason string `json:"reason" binding:"required"`
	}
	var data blockRequest
	if err := c.BindJSON(&data); err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}
	if err := u.service.BlockUser(id, userSessionIsAdmin, data.Reason); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (u *User) UnblockUser(c *gin.Context) {
	userSessionIsAdmin := c.GetBool("session_user_admin")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		err = app_errors.NewAppError(http.StatusBadRequest, "Invalid data in request", err)
		_ = c.Error(err)
		return
	}
	if err := u.service.UnblockUser(id, userSessionIsAdmin); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
