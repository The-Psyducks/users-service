package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"encoding/json"
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
	fmt.Printf("recorder body: %s ", recorder.Body.String())
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

func putUserPersonalInfo(router *router.Router, id string, user UserPersonalInfo) (int, error) {
	endpoint := fmt.Sprintf("/users/register/%s/personal-info", id)
	marshalledInfo, err := json.Marshal(user)

	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(marshalledInfo))

	if err != nil {
		return 0, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	return recorder.Code, nil
}

func putInterests(router *router.Router, id string, interests []int) (int, error) {
	payload := map[string][]int{
		"interests": interests,
	}
	marshalledInfo, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	endpoint := fmt.Sprintf("/users/register/%s/interests", id)
	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(marshalledInfo))

	if err != nil {
		return 0, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	return recorder.Code, nil
}

func completeRegistry(router *router.Router, id string) (int, UserProfile, error) {
	endpoint := fmt.Sprintf("/users/register/%s/complete", id)
	req, err := http.NewRequest("POST", endpoint, &bytes.Reader{})
	if err != nil {
		return 0, UserProfile{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := UserProfile{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return 0, UserProfile{}, err
	}

	return recorder.Code, result, nil
}

func CreateValidUser(router *router.Router, email string, personalInfo UserPersonalInfo, interests []int) (UserProfile, error) {
	res, err := getUserRegistryForSignUp(router, email)
	if err != nil {
		return UserProfile{}, err
	}

	err = sendEmailVerificationAndVerificateIt(router, res.Metadata.RegistrationId)
	if err != nil {
		return UserProfile{}, err
	}

	code, err := putUserPersonalInfo(router, res.Metadata.RegistrationId, personalInfo)
	if err != nil {
		return UserProfile{}, err
	}
	if code != http.StatusNoContent {
		return UserProfile{}, fmt.Errorf("error, status code adding personal info was %d", code)
	}

	code, err = putInterests(router, res.Metadata.RegistrationId, interests)
	if err != nil {
		return UserProfile{}, err
	}
	if code != http.StatusNoContent {
		return UserProfile{}, fmt.Errorf("error, status code adding interests was %d", code)
	}

	code, result, err := completeRegistry(router, res.Metadata.RegistrationId)
	if err != nil {
		return UserProfile{}, err
	}
	if code != http.StatusOK {
		return UserProfile{}, fmt.Errorf("error, status code completing registry was %d", code)
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

	code, err := putUserPersonalInfo(router, res.Metadata.RegistrationId, personalInfo)
	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}
	if code != http.StatusNoContent {
		return 0, ValidationErrorResponse{}, fmt.Errorf("error, status code adding personal info was %d", code)
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

func CreateExistingUser(router *router.Router, user UserPersonalInfo) (int, ErrorResponse, error) {
	marshalledInfo, err := json.Marshal(user)

	if err != nil {
		return 0, ErrorResponse{}, err
	}

	req, err := http.NewRequest("POST", "/users/register", bytes.NewReader(marshalledInfo))

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

func LoginValidUser(router *router.Router, loginReq LoginRequest) (int, LoginResponse, error) {
	marshalledInfo, err := json.Marshal(loginReq)

	if err != nil {
		return 0, LoginResponse{}, err
	}

	req, err := http.NewRequest("POST", "/users/login", bytes.NewReader(marshalledInfo))

	if err != nil {
		return 0, LoginResponse{}, err
	}

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := LoginResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return 0, LoginResponse{}, err
	}

	return recorder.Code, result, nil
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

func createAndLoginUser(router *router.Router, email string, user UserPersonalInfo, interestsIds []int) (int, LoginResponse, error) {
	_, err := CreateValidUser(router, email, user, interestsIds)
	if err != nil {
		return 0, LoginResponse{}, err
	}

	login := LoginRequest{
		Email:    email,
		Password: user.Password,
	}

	return LoginValidUser(router, login)
}

func getExistingUser(router *router.Router, username string, token string) (int, UserProfile, error) {
	req, err := http.NewRequest("GET", "/users/"+username, &bytes.Reader{})
	if err != nil {
		return 0, UserProfile{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := UserProfile{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return 0, UserProfile{}, err
	}
	return recorder.Code, result, nil
}

func getNotExistingUser(router *router.Router, username string, token string) (int, ErrorResponse, error) {
	req, err := http.NewRequest("GET", "/users/"+username, &bytes.Reader{})
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

func getRegisterOptions(router *router.Router) (int, RegisterOptions, error) {
	req, err := http.NewRequest("GET", "/users/register/locations", &bytes.Reader{})
	if err != nil {
		return 0, RegisterOptions{}, err
	}

	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	var locations  struct {
		Locations []Location `json:"locations"`
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &locations)
	if err != nil {
		return 0, RegisterOptions{}, err
	}
	
	req, err = http.NewRequest("GET", "/users/register/interests", &bytes.Reader{})
	if err != nil {
		return 0, RegisterOptions{}, err
	}
	
	recorder = httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	var interests  struct {
		Interests []Interest `json:"interests"`
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &interests)
	if err != nil {
		return 0, RegisterOptions{}, err
	}

	return recorder.Code, RegisterOptions{Locations: locations.Locations, Interests: interests.Interests}, nil
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

func AssertUserPrivateProfileIsUser(t *testing.T, email string, user UserPersonalInfo, location string, interests []string, profile UserProfile) {
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

func assertRegisterInstancePattern(t *testing.T, finalUrl string, expected string) {
	instancePattern := fmt.Sprintf(`^\/users\/register\/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\/%s$`, regexp.QuoteMeta(finalUrl))
	matched, err := regexp.MatchString(instancePattern, expected)
	assert.Equal(t, err, nil)
	assert.Equal(t, matched, true)
}
