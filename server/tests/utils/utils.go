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
	"users-service/src/router"
	"users-service/tests/constants"
	"users-service/tests/models"

	"github.com/go-playground/assert/v2"
)

func GetUserRegistryForSignUp(router *router.Router, email string) (models.ResolverSignUpResponse, error) {
	payload := map[string]string{
		"email": email,
	}
	marshalledInfo, err := json.Marshal(payload)
	if err != nil {
		return models.ResolverSignUpResponse{}, err
	}
	req, err := http.NewRequest("POST", "/users/resolver", bytes.NewReader(marshalledInfo))
	if err != nil {
		return models.ResolverSignUpResponse{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	res := models.ResolverSignUpResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	if err != nil {
		return models.ResolverSignUpResponse{}, err
	}

	if res.NextAuthStep != constants.SignUpAuthStep {
		return models.ResolverSignUpResponse{}, fmt.Errorf("error, next auth step was %s when it had to be %s", res.NextAuthStep, constants.SignUpAuthStep)
	}
	if res.Metadata.RegistrationId == "" {
		return models.ResolverSignUpResponse{}, fmt.Errorf("error, registration id was empty")
	}
	return res, nil
}

func SendEmailVerificationAndVerificateIt(router *router.Router, id string) error {
	endpoint := fmt.Sprintf("/users/register/%s/send-email", id)
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		return fmt.Errorf("error, status code sending email verification was %d, it had to be %d", recorder.Code, http.StatusNoContent)
	}

	payload := map[string]string{
		"pin": "421311",
	}
	marshalledInfo, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	endpoint = fmt.Sprintf("/users/register/%s/verify-email", id)
	req, err = http.NewRequest("POST", endpoint, bytes.NewReader(marshalledInfo))
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	recorder = httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		return fmt.Errorf("error, status code verifying email was %d", recorder.Code)
	}
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

	if recorder.Code != http.StatusNoContent {
		return fmt.Errorf("error, status code adding personal info was %d, expeted: %d", recorder.Code, http.StatusNoContent)
	}
	return nil
}

func PutValidInterests(router *router.Router, id string, interests []int) error {
	payload := map[string][]int{
		"interests": interests,
	}
	marshalledInfo, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/users/register/%s/interests", id)
	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(marshalledInfo))

	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		return fmt.Errorf("error, status code adding interests was %d, expected: %d", recorder.Code, http.StatusNoContent)
	}
	return nil
}

func CompleteValidRegistry(router *router.Router, id string) (models.UserPrivateProfile, error) {
	endpoint := fmt.Sprintf("/users/register/%s/complete", id)
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.UserPrivateProfile{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return models.UserPrivateProfile{}, err
	}

	if recorder.Code != http.StatusOK {
		return models.UserPrivateProfile{}, fmt.Errorf("error, status code completing registry was %d, expected: %d", recorder.Code, http.StatusOK)
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
	marshalledInfo, err := json.Marshal(personalInfo)
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(marshalledInfo))
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

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
	marshalledInfo, err := json.Marshal(payload)
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

	endpoint := fmt.Sprintf("/users/register/%s/interests", res.Metadata.RegistrationId)
	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(marshalledInfo))
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

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
	marshalledInfo, err := json.Marshal(loginReq)

	if err != nil {
		return models.LoginResponse{}, err
	}

	req, err := http.NewRequest("POST", "/users/login", bytes.NewReader(marshalledInfo))

	if err != nil {
		return models.LoginResponse{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.LoginResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return models.LoginResponse{}, err
	}

	if recorder.Code != http.StatusOK {
		return models.LoginResponse{}, fmt.Errorf("error, status code login was %d, expected: %d", recorder.Code, http.StatusOK)
	}

	return result, nil
}

func LoginInvalidUser(router *router.Router, loginReq models.LoginRequest) (int, models.ErrorResponse, error) {
	marshalledInfo, err := json.Marshal(loginReq)

	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	req, err := http.NewRequest("POST", "/users/login", bytes.NewReader(marshalledInfo))

	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

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
	req, err := http.NewRequest("GET", "/users/"+id, nil)
	if err != nil {
		return models.UserProfileResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.UserProfileResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return models.UserProfileResponse{}, err
	}

	if recorder.Code != http.StatusOK {
		return models.UserProfileResponse{}, fmt.Errorf("error, status code getting user was %d, expected: %d", recorder.Code, http.StatusOK)
	}
	return result, nil
}

func GetOwnProfile(router *router.Router, id string, token string) (models.UserPrivateProfile, error) {
	result, err := GetValidUser(router, id, token)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}
	if !result.OwnProfile {
		return models.UserPrivateProfile{}, fmt.Errorf("error, own profile was false")
	}
	profile := models.UserPrivateProfile{}
	jsonData, err := json.Marshal(result.Profile)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}
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
	if !result.OwnProfile {
		return models.UserPublicProfile{}, fmt.Errorf("error, own profile was false")
	}
	profile := models.UserPublicProfile{}
	jsonData, err := json.Marshal(result.Profile)
	if err != nil {
		return models.UserPublicProfile{}, err
	}
	err = json.Unmarshal(jsonData, &profile)
	if err != nil {
		return models.UserPublicProfile{}, err
	}
	return profile, nil
}

func GetNotExistingUser(router *router.Router, id string, token string) (models.ErrorResponse, error) {
	req, err := http.NewRequest("GET", "/users/"+id, nil)
	if err != nil {
		return models.ErrorResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return models.ErrorResponse{}, err
	}

	if recorder.Code != http.StatusNotFound {
		return models.ErrorResponse{}, fmt.Errorf("error, status code getting user was %d, expected: %d", recorder.Code, http.StatusNotFound)
	}
	return result, nil
}

func GetRegisterOptions(router *router.Router) (models.RegisterOptions, error) {
	req, err := http.NewRequest("GET", "/users/info/locations", nil)
	if err != nil {
		return models.RegisterOptions{}, err
	}

	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	var locations struct {
		Locations []models.Location `json:"locations"`
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &locations)
	if err != nil {
		return models.RegisterOptions{}, err
	}

	req, err = http.NewRequest("GET", "/users/info/interests", nil)
	if err != nil {
		return models.RegisterOptions{}, err
	}

	recorder = httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	var interests struct {
		Interests []models.Interest `json:"interests"`
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &interests)
	if err != nil {
		return models.RegisterOptions{}, err
	}

	if recorder.Code != http.StatusOK {
		return models.RegisterOptions{}, fmt.Errorf("error, status code getting register options was %d, expected: %d", recorder.Code, http.StatusOK)
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
	userInfo, err := json.Marshal(user)
	if err != nil {
		return models.UserPrivateProfile{}, fmt.Errorf("error, marshalling user info: %s", err.Error())
	}
	req, err := http.NewRequest("PUT", "/users/profile", bytes.NewReader(userInfo))
	if err != nil {
		return models.UserPrivateProfile{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.UserPrivateProfile{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return models.UserPrivateProfile{}, err
	}

	if recorder.Code != http.StatusOK {
		return models.UserPrivateProfile{}, fmt.Errorf("error, status code editing user profile was %d, expected: %d", recorder.Code, http.StatusNoContent)
	}

	return result, nil
}

func EditInvalidUserProfile(router *router.Router, token string, user models.EditUserProfileRequest) (int, models.ValidationErrorResponse, error) {
	userInfo, err := json.Marshal(user)
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}
	req, err := http.NewRequest("PUT", "/users/profile", bytes.NewReader(userInfo))
	if err != nil {
		return 0, models.ValidationErrorResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
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

func FollowValidUser(router *router.Router, id string, token string) error {
	req, err := http.NewRequest("POST", "/users/" + id + "/follow", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		return fmt.Errorf("error, status code following user was %d, expected: %d", recorder.Code, http.StatusNoContent)
	}
	return nil
}

func FollowInvalidUser(router *router.Router, id string, token string) (int, models.ErrorResponse, error) {
	req, err := http.NewRequest("POST", "/users/" + id + "/follow", nil)
	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func UnfollowValidUser(router *router.Router, id string, token string) error {
	req, err := http.NewRequest("DELETE", "/users/" + id + "/follow", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		return fmt.Errorf("error, status code unfollowing user was %d, expected: %d", recorder.Code, http.StatusNoContent)
	}
	return nil
}

func UnfollowInvalidUser(router *router.Router, id string, token string) (int, models.ErrorResponse, error) {
	req, err := http.NewRequest("DELETE", "/users/" + id + "/follow", nil)
	if err != nil {
		return 0, models.ErrorResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
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
		url := fmt.Sprintf("/users/%s/followers?time=%s&skip=%d&limit=%d", id, timestamp, skip, 20)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Add("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		router.Engine.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			return fmt.Errorf("error, status code getting followers was %d, expected: %d", recorder.Code, http.StatusOK)
		}

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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.ErrorResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := models.ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return models.ErrorResponse{}, err
	}
	if recorder.Code != http.StatusForbidden {
		return models.ErrorResponse{}, fmt.Errorf("error, status code getting followers was %d, expected: %d", recorder.Code, http.StatusOK)
	}
	return result, nil
}

func GetFollowing(router *router.Router, id string, token string) ([]models.FollowUserProfile, error) {
	var result []models.FollowUserProfile
	var currPagination models.Pagination

	fetchFollowing := func(skip int) error {
		timestamp := time.Now().UTC().Format(time.RFC3339Nano)
		url := fmt.Sprintf("/users/%s/following?time=%s&skip=%d&limit=%d", id, timestamp, skip, 20)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Add("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		router.Engine.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			return fmt.Errorf("error, status code getting following was %d, expected: %d", recorder.Code, http.StatusOK)
		}

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

func SearchUsers(router *router.Router, text string, token string, limit int) ([]models.FollowUserProfile, error) {
	var result []models.FollowUserProfile
	var currPagination models.Pagination
	
	fetchUsers := func(skip int) error {
		timestamp := time.Unix(time.Now().Unix()+1, 0).UTC().Format(time.RFC3339Nano)
		url := fmt.Sprintf("/users/search?text=%s&time=%s&skip=%d&limit=%d",text, timestamp, skip, limit)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
	
		req.Header.Add("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		router.Engine.ServeHTTP(recorder, req)
	
		if recorder.Code != http.StatusOK {
			return fmt.Errorf("error, status code getting following was %d, expected: %d", recorder.Code, http.StatusOK)
		}
	
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
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
	
		req.Header.Add("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		router.Engine.ServeHTTP(recorder, req)
	
		if recorder.Code != http.StatusOK {
			return fmt.Errorf("error, status code getting following was %d, expected: %d", recorder.Code, http.StatusOK)
		}
	
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
