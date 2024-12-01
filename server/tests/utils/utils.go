package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
	"users-service/src/auth"
	"users-service/src/router"
	"users-service/tests/models"

	"github.com/go-playground/assert/v2"
)

func GetUserRegistryForSignUp(router *router.Router, email string) (models.ResolverSignUpResponse, error) {
	payload := map[string]string{
		"email": email,
	}
	marshalledInfo, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/users/resolver", bytes.NewReader(marshalledInfo))

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	res := models.ResolverSignUpResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &res)
	if err != nil {
		return models.ResolverSignUpResponse{}, err
	}

	return res, nil
}

func SendEmailVerificationAndVerificateIt(router *router.Router, id string) error {
	endpoint := fmt.Sprintf("/users/register/%s/send-email", id)
	req, _ := http.NewRequest("POST", endpoint, nil)

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	payload := map[string]string{
		"pin": "421311",
	}
	marshalledInfo, _ := json.Marshal(payload)

	endpoint = fmt.Sprintf("/users/register/%s/verify-email", id)
	req, _ = http.NewRequest("POST", endpoint, bytes.NewReader(marshalledInfo))

	req.Header.Add("content-type", "application/json")
	recorder = httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	return nil
}

func PutValidUserPersonalInfo(router *router.Router, id string, user models.UserPersonalInfo) error {
	endpoint := fmt.Sprintf("/users/register/%s/personal-info", id)
	marshalledInfo, err := json.Marshal(user)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(marshalledInfo))

	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	return nil
}

func PutValidInterests(router *router.Router, id string, interests []int) error {
	payload := map[string][]int{
		"interests": interests,
	}
	marshalledInfo, _ := json.Marshal(payload)

	endpoint := fmt.Sprintf("/users/register/%s/interests", id)
	req, _ := http.NewRequest("PUT", endpoint, bytes.NewReader(marshalledInfo))

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	return nil
}

func CompleteValidRegistry(router *router.Router, id string) (models.UserPrivateProfile, error) {
	endpoint := fmt.Sprintf("/users/register/%s/complete", id)
	req, _ := http.NewRequest("POST", endpoint, nil)

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.UserPrivateProfile{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return models.UserPrivateProfile{}, err
	}

	return result, nil
}

func CreateValidUser(router *router.Router, email string, personalInfo models.UserPersonalInfo, interests []int) (models.UserPrivateProfile, error) {
	res, err := GetUserRegistryForSignUp(router, email)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}

	err = SendEmailVerificationAndVerificateIt(router, res.Metadata.RegistrationId)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}

	err = PutValidUserPersonalInfo(router, res.Metadata.RegistrationId, personalInfo)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}
	err = PutValidInterests(router, res.Metadata.RegistrationId, interests)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}
	result, err := CompleteValidRegistry(router, res.Metadata.RegistrationId)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}

	return result, nil
}

func CreateUserWithInvalidPersonalInfo(router *router.Router, email string, personalInfo models.UserPersonalInfo) (int, models.ValidationErrorResponse, error) {
	res, err := GetUserRegistryForSignUp(router, email)
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

	err = SendEmailVerificationAndVerificateIt(router, res.Metadata.RegistrationId)
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

	endpoint := fmt.Sprintf("/users/register/%s/personal-info", res.Metadata.RegistrationId)
	marshalledInfo, _ := json.Marshal(personalInfo)

	req, _ := http.NewRequest("PUT", endpoint, bytes.NewReader(marshalledInfo))

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ValidationErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func CreateUserWithInvalidInterests(router *router.Router, email string, personalInfo models.UserPersonalInfo, interestsIds []int) (int, models.ValidationErrorResponse, error) {
	res, err := GetUserRegistryForSignUp(router, email)
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

	err = SendEmailVerificationAndVerificateIt(router, res.Metadata.RegistrationId)
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

	err = PutValidUserPersonalInfo(router, res.Metadata.RegistrationId, personalInfo)
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

	payload := map[string][]int{
		"interests": interestsIds,
	}
	marshalledInfo, _ := json.Marshal(payload)

	endpoint := fmt.Sprintf("/users/register/%s/interests", res.Metadata.RegistrationId)
	req, _ := http.NewRequest("PUT", endpoint, bytes.NewReader(marshalledInfo))

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ValidationErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func LoginValidUser(router *router.Router, loginReq models.LoginRequest) (models.LoginResponse, error) {
	marshalledInfo, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewReader(marshalledInfo))

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.LoginResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return models.LoginResponse{}, err
	}

	return result, nil
}

func LoginInvalidUser(router *router.Router, loginReq models.LoginRequest) (int, models.ErrorResponse, error) {
	marshalledInfo, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest("POST", "/users/login", bytes.NewReader(marshalledInfo))

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func CreateAndLoginUser(router *router.Router, email string, user models.UserPersonalInfo, interestsIds []int) (models.LoginResponse, error) {
	_, err := CreateValidUser(router, email, user, interestsIds)
	if err != nil {
		return models.LoginResponse{}, err
	}

	login := models.LoginRequest{
		Email:    email,
		Password: user.Password,
	}

	return LoginValidUser(router, login)
}

func GetValidUser(router *router.Router, id string, token string) (models.UserProfileResponse, error) {
	req, _ := http.NewRequest("GET", "/users/"+id, nil)

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.UserProfileResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return models.UserProfileResponse{}, err
	}

	return result, nil
}

func GetValidUserInformation(router *router.Router, id string, adminToken string) (models.UserInformationResponse, error) {
	req, _ := http.NewRequest("GET", "/users/"+id+"/information", nil)

	req.Header.Add("Authorization", "Bearer "+adminToken)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.UserInformationResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return models.UserInformationResponse{}, err
	}

	return result, nil
}

func GetInvalidUserInformation(router *router.Router, id string, token string) (int, models.ErrorResponse, error) {
	req, _ := http.NewRequest("GET", "/users/"+id+"/information", nil)
	
	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func GetOwnProfile(router *router.Router, id string, token string) (models.UserPrivateProfile, error) {
	result, err := GetValidUser(router, id, token)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}

	profile := models.UserPrivateProfile{}
	jsonData, _ := json.Marshal(result.Profile)
	
	err = json.Unmarshal(jsonData, &profile)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}
	return profile, nil
}

func GetAnotherUserProfile(router *router.Router, id string, token string) (models.UserPublicProfile, error) {
	result, err := GetValidUser(router, id, token)
	if err != nil {
		return models.UserPublicProfile{}, err
	}

	profile := models.UserPublicProfile{}
	jsonData, _ := json.Marshal(result.Profile)
	
	err = json.Unmarshal(jsonData, &profile)
	if err != nil {
		return models.UserPublicProfile{}, err
	}
	return profile, nil
}

func GetNotExistingUser(router *router.Router, id string, token string) (models.ErrorResponse, error) {
	req, _ := http.NewRequest("GET", "/users/"+id, nil)
	
	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return models.ErrorResponse{}, err
	}

	return result, nil
}

func GetRegisterOptions(router *router.Router) (models.RegisterOptions, error) {
	req, _ := http.NewRequest("GET", "/users/info/locations", nil)
	
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	var locations struct {
		Locations []models.Location `json:"locations"`
	}
	err := json.Unmarshal(recorder.Body.Bytes(), &locations)
	if err != nil {
		return models.RegisterOptions{}, err
	}

	req, _ = http.NewRequest("GET", "/users/info/interests", nil)
	
	recorder = httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	var interests struct {
		Interests []models.Interest `json:"interests"`
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &interests)
	if err != nil {
		return models.RegisterOptions{}, err
	}

	return models.RegisterOptions{Locations: locations.Locations, Interests: interests.Interests}, nil
}

func GetLocationAndInterestsNames(registerOptions models.RegisterOptions, locationId int, interestsIds []int) (string, []string) {
	var location string
	interests := make([]string, len(interestsIds))

	for _, loc := range registerOptions.Locations {
		if loc.Id == locationId {
			location = loc.Name
			break
		}
	}

	for i, interestId := range interestsIds {
		for _, interest := range registerOptions.Interests {
			if interest.Id == interestId {
				interests[i] = interest.Name
				break
			}
		}
	}
	return location, interests
}

func EditValidUserProfile(router *router.Router, token string, user models.EditUserProfileRequest) (models.UserPrivateProfile, error) {
	userInfo, _ := json.Marshal(user)

	req, _ := http.NewRequest("PUT", "/users/profile", bytes.NewReader(userInfo))
	
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.UserPrivateProfile{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}

	return result, nil
}

func EditInvalidUserProfile(router *router.Router, token string, user models.EditUserProfileRequest) (int, models.ValidationErrorResponse, error) {
	userInfo, _ := json.Marshal(user)

	req, _ := http.NewRequest("PUT", "/users/profile", bytes.NewReader(userInfo))
	
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ValidationErrorResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func FollowValidUser(router *router.Router, id string, token string) error {
	req, _ := http.NewRequest("POST", "/users/" + id + "/follow", nil)
	
	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	return nil
}

func FollowInvalidUser(router *router.Router, id string, token string) (int, models.ErrorResponse, error) {
	req, _ := http.NewRequest("POST", "/users/" + id + "/follow", nil)
	
	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func UnfollowValidUser(router *router.Router, id string, token string) error {
	req, _ := http.NewRequest("DELETE", "/users/" + id + "/follow", nil)

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	return nil
}

func UnfollowInvalidUser(router *router.Router, id string, token string) (int, models.ErrorResponse, error) {
	req, _ := http.NewRequest("DELETE", "/users/" + id + "/follow", nil)
	
	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func GetFollowers(router *router.Router, id string, token string) ([]models.FollowUserProfile, error) {
	var result []models.FollowUserProfile
	var currPagination models.Pagination

	fetchFollowers := func(skip int) error {
		timestamp := time.Now().UTC().Format(time.RFC3339Nano)
		url := fmt.Sprintf("/users/%s/followers?time=%s&skip=%d&limit=%d", id, timestamp, skip, 1)
		req, _ := http.NewRequest("GET", url, nil)
		
		req.Header.Add("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		router.Engine.ServeHTTP(recorder, req)

		newResult := models.FollowersResponse{}
		if err := json.Unmarshal(recorder.Body.Bytes(), &newResult); err != nil {
			return err
		}
		result = append(result, newResult.Followers...)
		currPagination = newResult.Pagination
		return nil
	}

	if err := fetchFollowers(0); err != nil {
		return nil, err
	}

	for currPagination.NextOffset != 0 {
		if err := fetchFollowers(currPagination.NextOffset); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func GetFollowersForInvalidUser(router *router.Router, id string, token string) (models.ErrorResponse, error) {
	timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	url := fmt.Sprintf("/users/%s/following?timestamp=%s&skip=%d&limit=%d", id, timestamp, 0, 20)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return models.ErrorResponse{}, err
	}

	return result, nil
}

func GetFollowing(router *router.Router, id string, token string) ([]models.FollowUserProfile, error) {
	var result []models.FollowUserProfile
	var currPagination models.Pagination

	fetchFollowing := func(skip int) error {
		timestamp := time.Now().UTC().Format(time.RFC3339Nano)
		url := fmt.Sprintf("/users/%s/following?time=%s&skip=%d&limit=%d", id, timestamp, skip, 20)
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		router.Engine.ServeHTTP(recorder, req)

		newResult := models.FollowingResponse{}
		if err := json.Unmarshal(recorder.Body.Bytes(), &newResult); err != nil {
			return err
		}

		result = append(result, newResult.Following...)
		currPagination = newResult.Pagination
		return nil
	}

	if err := fetchFollowing(0); err != nil {
		return nil, err
	}

	for currPagination.NextOffset != 0 {
		if err := fetchFollowing(currPagination.NextOffset); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func GetFollowingForInvalidUser(router *router.Router, id string, token string) (models.ErrorResponse, error) {
	timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	url := fmt.Sprintf("/users/%s/following?time=%s&skip=%d&limit=%d", id, timestamp, 0, 20)
	req, _ := http.NewRequest("GET", url, nil)
	
	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return models.ErrorResponse{}, err
	}
	return result, nil
}


func GetAmountOfFollowersInTimeRange(router *router.Router, id, token, startTime, endTime string) (int, error) {
	url := fmt.Sprintf("/users/metrics/%s/followers?time=%s&end_time=%s", id, startTime, endTime)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	var result struct {
		Amount int `json:"new_followers"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		return 0, err
	}
	return result.Amount, nil
}

func GetAmountOfFollowersInTimeRangeInvalid(router *router.Router, id, token, startTime, endTime string) (int, error) {
	url := fmt.Sprintf("/users/metrics/%s/followers?time=%s&end_time=%s", id, startTime, endTime)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	var result models.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		return 0, err
	}
	return result.Status, nil
}

func SearchUsers(router *router.Router, text string, token string, limit int) ([]models.FollowUserProfile, error) {
	var result []models.FollowUserProfile
	var currPagination models.Pagination
	
	fetchUsers := func(skip int) error {
		timestamp := time.Unix(time.Now().Unix()+1, 0).UTC().Format(time.RFC3339Nano)
		url := fmt.Sprintf("/users/search?text=%s&time=%s&skip=%d&limit=%d",text, timestamp, skip, limit)
		req, _ := http.NewRequest("GET", url, nil)
	
		req.Header.Add("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		router.Engine.ServeHTTP(recorder, req)
	
		newResult := models.PaginationResponse[models.FollowUserProfile]{}
		if err := json.Unmarshal(recorder.Body.Bytes(), &newResult); err != nil {
			return err
		}
	
		result = append(result, newResult.Data...)
		currPagination = newResult.Pagination
		return nil
	}
	
	if err := fetchUsers(0); err != nil {
		return nil, err
	}
	
	for currPagination.NextOffset != 0 {
		if err := fetchUsers(currPagination.NextOffset); err != nil {
			return nil, err
		}
	}
	
	return result, nil
}

func GetAllUserRecommendations(router *router.Router, token string, limit int) ([]models.FollowUserProfile, error) {
	var result []models.FollowUserProfile
	var currPagination models.Pagination
	
	fetchUsers := func(skip int) error {
		timestamp := time.Unix(time.Now().Unix()+1, 0).UTC().Format(time.RFC3339Nano)
		url := fmt.Sprintf("/users/recommendations?&time=%s&skip=%d&limit=%d", timestamp, skip, limit)
		req, _ := http.NewRequest("GET", url, nil)
	
		req.Header.Add("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		router.Engine.ServeHTTP(recorder, req)
	
		newResult := models.PaginationResponse[models.FollowUserProfile]{}
		if err := json.Unmarshal(recorder.Body.Bytes(), &newResult); err != nil {
			return err
		}
	
		result = append(result, newResult.Data...)
		currPagination = newResult.Pagination
		return nil
	}
	
	if err := fetchUsers(0); err != nil {
		return nil, err
	}
	
	for currPagination.NextOffset != 0 {
		if err := fetchUsers(currPagination.NextOffset); err != nil {
			return nil, err
		}
	}
	
	return result, nil
}

func GetAllUsers(router *router.Router, token string, limit int) ([]models.UserPublicProfile, error) {
	var result []models.UserPublicProfile
	var currPagination models.Pagination
	
	fetchUsers := func(skip int) error {
		timestamp := time.Unix(time.Now().Unix()+1, 0).UTC().Format(time.RFC3339Nano)
		url := fmt.Sprintf("/users/all?time=%s&skip=%d&limit=%d", timestamp, skip, limit)
		req, _ := http.NewRequest("GET", url, nil)
	
		req.Header.Add("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		router.Engine.ServeHTTP(recorder, req)
	
		newResult := models.PaginationResponse[models.UserPublicProfile]{}
		if err := json.Unmarshal(recorder.Body.Bytes(), &newResult); err != nil {
			return err
		}
	
		result = append(result, newResult.Data...)
		currPagination = newResult.Pagination
		return nil
	}

	if err := fetchUsers(0); err != nil {
		return nil, err
	}
	
	for currPagination.NextOffset != 0 {
		if err := fetchUsers(currPagination.NextOffset); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func GetAllUsersInvalidToken(router *router.Router, token string, limit int) (int, models.ErrorResponse, error) {
	timestamp := time.Unix(time.Now().Unix()+1, 0).UTC().Format(time.RFC3339Nano)
	url := fmt.Sprintf("/users/all?time=%s&skip=%d&limit=%d", timestamp, 0, limit)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func LoginAdmin() (string, error) {
	token, err := auth.GenerateToken("edf533b4-6ea5-414f-8442-320f60428b8e", true)
	return token, err
}

func BlockUser(router *router.Router, id string, reason string, token string) error {
	payload := map[string]string{
		"reason": reason,
	}
	marshalledInfo, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/users/"+id+"/block", bytes.NewReader(marshalledInfo))
	
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	return nil
}

func BlockInvalidUser(router *router.Router, id string, reason string, token string) (int, models.ErrorResponse, error) {
	var req *http.Request

	if reason != "" {
		payload := map[string]string{
			"reason": reason,
		}
		marshalledInfo, _ := json.Marshal(payload)
		req, _ = http.NewRequest("POST", "/users/"+id+"/block", bytes.NewReader(marshalledInfo))
	} else {
		req, _ = http.NewRequest("POST", "/users/"+id+"/block", nil)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func UnblockUser(router *router.Router, id string, token string) error {
	req, _ := http.NewRequest("POST", "/users/"+id+"/unblock", nil)

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	return nil
}

func UnblockInvalidUser(router *router.Router, id string, token string) (int, models.ErrorResponse, error) {
	req, _ := http.NewRequest("POST", "/users/"+id+"/unblock", nil)

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err := json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func AssertUserPrivateProfileIsUser(t *testing.T, email string, user models.UserPersonalInfo, location string, interests []string, profile models.UserPrivateProfile) {
	assert.Equal(t, user.FirstName, profile.FirstName)
	assert.Equal(t, user.LastName, profile.LastName)
	assert.Equal(t, user.UserName, profile.UserName)
	assert.Equal(t, email, profile.Email)
	assert.Equal(t, location, profile.Location)

	assert.Equal(t, len(profile.Interests), len(interests))
	for _, interest := range interests {
		found := false
		for _, profileInterest := range profile.Interests {
			if interest == profileInterest {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("models.Interest %s not found in profile interests", interest)
		}
	}
}

func AssertPrivateUsersAreEqual(t *testing.T, expected models.UserPrivateProfile, actual models.UserPrivateProfile) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.FirstName, actual.FirstName)
	assert.Equal(t, expected.LastName, actual.LastName)
	assert.Equal(t, expected.UserName, actual.UserName)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Location, actual.Location)
	assert.Equal(t, expected.Followers, actual.Followers)
	assert.Equal(t, expected.Following, actual.Following)

	assert.Equal(t, len(expected.Interests), len(actual.Interests))
	for _, interest := range expected.Interests {
		found := false
		for _, actualInterest := range actual.Interests {
			if interest == actualInterest {
				found = true
				break
			}
		}
		assert.Equal(t, found, true)
	}
}

func AssertUserPublicProfileIsUser(t *testing.T, user models.UserPersonalInfo, location string, profile models.UserPublicProfile) {
	assert.Equal(t, user.FirstName, profile.FirstName)
	assert.Equal(t, user.LastName, profile.LastName)
	assert.Equal(t, user.UserName, profile.UserName)
	assert.Equal(t, location, profile.Location)
}

func AssertInterestsNamesAreCorrectIds(t *testing.T, registerOptions models.RegisterOptions, interestsIds []int, interests []string) {
	assert.Equal(t, len(interestsIds), len(interests))
	for _, interestId := range interestsIds {
		found := false
		for _, interest := range registerOptions.Interests {
			if interest.Id == interestId {
				for _, name := range interests {
					if name == interest.Name {
						found = true
						break
					}
				}
				break
			}
		}
		assert.Equal(t, found, true)
	}
}

func AssertLocationNameIsCorrectId(t *testing.T, registerOptions models.RegisterOptions, locationId int, location string) {
	found := false
	for _, loc := range registerOptions.Locations {
		if loc.Id == locationId {
			found = true
			assert.Equal(t, loc.Name, location)
			break
		}
	}
	assert.Equal(t, found, true)
}

func AssertRegisterInstancePattern(t *testing.T, finalUrl string, expected string) {
	instancePattern := fmt.Sprintf(`^\/users\/register\/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\/%s$`, regexp.QuoteMeta(finalUrl))
	matched, err := regexp.MatchString(instancePattern, expected)
	assert.Equal(t, err, nil)
	assert.Equal(t, matched, true)
}

func AssertListsAreEqual(t *testing.T, expected []models.FollowUserProfile, actual []models.FollowUserProfile) {
	assert.Equal(t, len(expected), len(actual))

	for _, e := range expected {
		found := false
		for _, a := range actual {
			if e.Profile.Id.String() == a.Profile.Id.String() {
				found = true
				break
			}
		}
		assert.Equal(t, found, true)
	}
}

func PrivateUserToPublic(user models.UserPrivateProfile) models.UserPublicProfile {
	return models.UserPublicProfile{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserName:  user.UserName,
		Location:  user.Location,
		Followers: user.Followers,
		Following: user.Following,
	}
}
