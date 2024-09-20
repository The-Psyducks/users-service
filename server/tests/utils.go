package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"users-service/src/router"

	"github.com/go-playground/assert/v2"
)

func getUserRegistryForSignUp(router *router.Router, email string) (ResolverSignUpResponse, error) {
	payload := map[string]string{
		"email": email,
	}
	marshalledInfo, err := json.Marshal(payload)
	if err != nil {
		return ResolverSignUpResponse{}, err
	}
	req, err := http.NewRequest("POST", "/users/resolver", bytes.NewReader(marshalledInfo))
	if err != nil {
		return ResolverSignUpResponse{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	res := ResolverSignUpResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	if err != nil {
		return ResolverSignUpResponse{}, err
	}

	if res.NextAuthStep != SignUpAuthStep {
		return ResolverSignUpResponse{}, fmt.Errorf("error, next auth step was %s when it had to be %s", res.NextAuthStep, SignUpAuthStep)
	}
	if res.Metadata.RegistrationId == "" {
		return ResolverSignUpResponse{}, fmt.Errorf("error, registration id was empty")
	}
	return res, nil
}

func sendEmailVerificationAndVerificateIt(router *router.Router, id string) error {
	endpoint := fmt.Sprintf("/users/register/%s/send-email", id)
	req, err := http.NewRequest("POST", endpoint, &bytes.Reader{})
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
		"pin": "123456",
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

func putValidUserPersonalInfo(router *router.Router, id string, user UserPersonalInfo) error {
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

func putValidInterests(router *router.Router, id string, interests []int) error {
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

func completeValidRegistry(router *router.Router, id string) (UserPrivateProfile, error) {
	endpoint := fmt.Sprintf("/users/register/%s/complete", id)
	req, err := http.NewRequest("POST", endpoint, &bytes.Reader{})
	if err != nil {
		return UserPrivateProfile{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := UserPrivateProfile{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return UserPrivateProfile{}, err
	}

	if recorder.Code != http.StatusOK {
		return UserPrivateProfile{}, fmt.Errorf("error, status code completing registry was %d, expected: %d", recorder.Code, http.StatusOK)
	}
	return result, nil
}

func CreateValidUser(router *router.Router, email string, personalInfo UserPersonalInfo, interests []int) (UserPrivateProfile, error) {
	res, err := getUserRegistryForSignUp(router, email)
	if err != nil {
		return UserPrivateProfile{}, err
	}

	err = sendEmailVerificationAndVerificateIt(router, res.Metadata.RegistrationId)
	if err != nil {
		return UserPrivateProfile{}, err
	}

	err = putValidUserPersonalInfo(router, res.Metadata.RegistrationId, personalInfo)
	if err != nil {
		return UserPrivateProfile{}, err
	}
	err = putValidInterests(router, res.Metadata.RegistrationId, interests)
	if err != nil {
		return UserPrivateProfile{}, err
	}
	result, err := completeValidRegistry(router, res.Metadata.RegistrationId)
	if err != nil {
		return UserPrivateProfile{}, err
	}

	return result, nil
}

func CreateUserWithInvalidPersonalInfo(router *router.Router, email string, personalInfo UserPersonalInfo) (int, ValidationErrorResponse, error) {
	res, err := getUserRegistryForSignUp(router, email)
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	err = sendEmailVerificationAndVerificateIt(router, res.Metadata.RegistrationId)
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	endpoint := fmt.Sprintf("/users/register/%s/personal-info", res.Metadata.RegistrationId)
	marshalledInfo, err := json.Marshal(personalInfo)
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(marshalledInfo))
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := ValidationErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func CreateUserWithInvalidInterests(router *router.Router, email string, personalInfo UserPersonalInfo, interestsIds []int) (int, ValidationErrorResponse, error) {
	res, err := getUserRegistryForSignUp(router, email)
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	err = sendEmailVerificationAndVerificateIt(router, res.Metadata.RegistrationId)
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	err = putValidUserPersonalInfo(router, res.Metadata.RegistrationId, personalInfo)
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	payload := map[string][]int{
		"interests": interestsIds,
	}
	marshalledInfo, err := json.Marshal(payload)
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	endpoint := fmt.Sprintf("/users/register/%s/interests", res.Metadata.RegistrationId)
	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(marshalledInfo))
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := ValidationErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func LoginValidUser(router *router.Router, loginReq LoginRequest) (LoginResponse, error) {
	marshalledInfo, err := json.Marshal(loginReq)

	if err != nil {
		return LoginResponse{}, err
	}

	req, err := http.NewRequest("POST", "/users/login", bytes.NewReader(marshalledInfo))

	if err != nil {
		return LoginResponse{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := LoginResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return LoginResponse{}, err
	}

	if recorder.Code != http.StatusOK {
		return LoginResponse{}, fmt.Errorf("error, status code login was %d, expected: %d", recorder.Code, http.StatusOK)
	}

	return result, nil
}

func LoginInvalidUser(router *router.Router, loginReq LoginRequest) (int, ErrorResponse, error) {
	marshalledInfo, err := json.Marshal(loginReq)

	if err != nil {
		return 0, ErrorResponse{}, err
	}

	req, err := http.NewRequest("POST", "/users/login", bytes.NewReader(marshalledInfo))

	if err != nil {
		return 0, ErrorResponse{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return 0, ErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func createAndLoginUser(router *router.Router, email string, user UserPersonalInfo, interestsIds []int) (LoginResponse, error) {
	_, err := CreateValidUser(router, email, user, interestsIds)
	if err != nil {
		return LoginResponse{}, err
	}

	login := LoginRequest{
		Email:    email,
		Password: user.Password,
	}

	return LoginValidUser(router, login)
}

func getValidUser(router *router.Router, username string, token string) (UserProfileResponse, error) {
	req, err := http.NewRequest("GET", "/users/"+username, &bytes.Reader{})
	if err != nil {
		return UserProfileResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := UserProfileResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return UserProfileResponse{}, err
	}

	if recorder.Code != http.StatusOK {
		return UserProfileResponse{}, fmt.Errorf("error, status code getting user was %d, expected: %d", recorder.Code, http.StatusOK)
	}
	return result, nil
}

func getOwnProfile(router *router.Router, username string, token string) (UserPrivateProfile, error) {
	result, err := getValidUser(router, username, token)
	if err != nil {
		return UserPrivateProfile{}, err
	}
	if !result.OwnProfile {
		return UserPrivateProfile{}, fmt.Errorf("error, own profile was false")
	}
	profile := UserPrivateProfile{}
	jsonData, err := json.Marshal(result.Profile)
	if err != nil {
		return UserPrivateProfile{}, err
	}
	err = json.Unmarshal(jsonData, &profile)
	if err != nil {
		return UserPrivateProfile{}, err
	}
	return profile, nil
}

func getAnotherUserProfile(router *router.Router, username string, token string) (UserPublicProfile, error) {
	result, err := getValidUser(router, username, token)
	if err != nil {
		return UserPublicProfile{}, err
	}
	if !result.OwnProfile {
		return UserPublicProfile{}, fmt.Errorf("error, own profile was false")
	}
	profile := UserPublicProfile{}
	jsonData, err := json.Marshal(result.Profile)
	if err != nil {
		return UserPublicProfile{}, err
	}
	err = json.Unmarshal(jsonData, &profile)
	if err != nil {
		return UserPublicProfile{}, err
	}
	return profile, nil
}

func getNotExistingUser(router *router.Router, username string, token string) (ErrorResponse, error) {
	req, err := http.NewRequest("GET", "/users/"+username, &bytes.Reader{})
	if err != nil {
		return ErrorResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return ErrorResponse{}, err
	}

	if recorder.Code != http.StatusNotFound {
		return ErrorResponse{}, fmt.Errorf("error, status code getting user was %d, expected: %d", recorder.Code, http.StatusNotFound)
	}
	return result, nil
}

func getRegisterOptions(router *router.Router) (RegisterOptions, error) {
	req, err := http.NewRequest("GET", "/users/register/locations", &bytes.Reader{})
	if err != nil {
		return RegisterOptions{}, err
	}

	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	var locations struct {
		Locations []Location `json:"locations"`
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &locations)
	if err != nil {
		return RegisterOptions{}, err
	}

	req, err = http.NewRequest("GET", "/users/register/interests", &bytes.Reader{})
	if err != nil {
		return RegisterOptions{}, err
	}

	recorder = httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	var interests struct {
		Interests []Interest `json:"interests"`
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &interests)
	if err != nil {
		return RegisterOptions{}, err
	}

	if recorder.Code != http.StatusOK {
		return RegisterOptions{}, fmt.Errorf("error, status code getting register options was %d, expected: %d", recorder.Code, http.StatusOK)
	}

	return RegisterOptions{Locations: locations.Locations, Interests: interests.Interests}, nil
}

func getLocationAndInterestsNames(registerOptions RegisterOptions, locationId int, interestsIds []int) (string, []string) {
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

func followValidUser(router *router.Router, username string, token string) error {
	payload := map[string]string{
		"username": username,
	}
	marshalledInfo, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "/users/follow", bytes.NewBuffer(marshalledInfo))
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

func followInvalidUser(router *router.Router, username string, token string) (int, ErrorResponse, error) {
	payload := map[string]string{
		"username": username,
	}
	marshalledInfo, err := json.Marshal(payload)
	if err != nil {
		return 0, ErrorResponse{}, err
	}
	req, err := http.NewRequest("POST", "/users/follow", bytes.NewBuffer(marshalledInfo))
	if err != nil {
		return 0, ErrorResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, ErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func unfollowValidUser(router *router.Router, username string, token string) error {
	payload := map[string]string{
		"username": username,
	}
	marshalledInfo, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", "/users/follow", bytes.NewBuffer(marshalledInfo))
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

func unfollowInvalidUser(router *router.Router, username string, token string) (int, ErrorResponse, error) {
	payload := map[string]string{
		"username": username,
	}
	marshalledInfo, err := json.Marshal(payload)
	if err != nil {
		return 0, ErrorResponse{}, err
	}
	req, err := http.NewRequest("DELETE", "/users/follow", bytes.NewBuffer(marshalledInfo))
	if err != nil {
		return 0, ErrorResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return 0, ErrorResponse{}, err
	}

	return recorder.Code, result, nil
}

func getFollowers(router *router.Router, username string, token string) (FollowersResponse, error) {
	req, err := http.NewRequest("GET", "/users/"+username+"/followers", &bytes.Reader{})
	if err != nil {
		return FollowersResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := FollowersResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return FollowersResponse{}, err
	}
	if recorder.Code != http.StatusOK {
		return FollowersResponse{}, fmt.Errorf("error, status code getting followers was %d, expected: %d", recorder.Code, http.StatusOK)
	}
	return result, nil
}

func getFollowersForInvalidUser(router *router.Router, username string, token string) (ErrorResponse, error) {
	req, err := http.NewRequest("GET", "/users/"+username+"/followers", &bytes.Reader{})
	if err != nil {
		return ErrorResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return ErrorResponse{}, err
	}
	if recorder.Code != http.StatusForbidden {
		return ErrorResponse{}, fmt.Errorf("error, status code getting followers was %d, expected: %d", recorder.Code, http.StatusOK)
	}
	return result, nil
}

func getFollowing(router *router.Router, username string, token string) (FollowingResponse, error) {
	req, err := http.NewRequest("GET", "/users/"+username+"/following", &bytes.Reader{})
	if err != nil {
		return FollowingResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := FollowingResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return FollowingResponse{}, err
	}
	if recorder.Code != http.StatusOK {
		return FollowingResponse{}, fmt.Errorf("error, status code getting followers was %d, expected: %d", recorder.Code, http.StatusOK)
	}
	return result, nil
}

func AssertUserPrivateProfileIsUser(t *testing.T, email string, user UserPersonalInfo, location string, interests []string, profile UserPrivateProfile) {
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
			t.Errorf("Interest %s not found in profile interests", interest)
		}
	}
}

func AssertUserPublicProfileIsUser(t *testing.T, user UserPersonalInfo, location string, profile UserPublicProfile) {
	assert.Equal(t, user.FirstName, profile.FirstName)
	assert.Equal(t, user.LastName, profile.LastName)
	assert.Equal(t, user.UserName, profile.UserName)
	assert.Equal(t, location, profile.Location)
}

func assertRegisterInstancePattern(t *testing.T, finalUrl string, expected string) {
	instancePattern := fmt.Sprintf(`^\/users\/register\/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\/%s$`, regexp.QuoteMeta(finalUrl))
	matched, err := regexp.MatchString(instancePattern, expected)
	assert.Equal(t, err, nil)
	assert.Equal(t, matched, true)
}

func assertListsAreEqual(t *testing.T, expected []FollowUserProfile, actual []FollowUserProfile) {
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

func privateUserToPublic(user UserPrivateProfile) UserPublicProfile {
	return UserPublicProfile{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserName:  user.UserName,
		Location:  user.Location,
		Followers: user.Followers,
		Following: user.Following,
	}
}
