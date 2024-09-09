package tests

import (
	"bytes"
	"testing"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"users-service/src/router"
	"github.com/go-playground/assert/v2"
)

func CreateValidUser(router *router.Router, user User) (int, UserProfile, error) {
	marshalledSnap, err := json.Marshal(user)

	if err != nil {
		return 0, UserProfile{}, err
	}

	req, err := http.NewRequest("POST", "/users/register", bytes.NewReader(marshalledSnap))

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

func CreateInvalidUser(router *router.Router, user User) (int, ValidationErrorResponse, error) {
	marshalledSnap, err := json.Marshal(user)

	if err != nil {
		return 0, ValidationErrorResponse{}, err
	}

	req, err := http.NewRequest("POST", "/users/register", bytes.NewReader(marshalledSnap))

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

func CreateExistingUser(router *router.Router, user User) (int, ErrorResponse, error) {
	marshalledSnap, err := json.Marshal(user)

	if err != nil {
		return 0, ErrorResponse{}, err
	}

	req, err := http.NewRequest("POST", "/users/register", bytes.NewReader(marshalledSnap))

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
	marshalledSnap, err := json.Marshal(loginReq)

	if err != nil {
		return 0, LoginResponse{}, err
	}

	req, err := http.NewRequest("POST", "/users/login", bytes.NewReader(marshalledSnap))

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
	marshalledSnap, err := json.Marshal(loginReq)

	if err != nil {
		return 0, ErrorResponse{}, err
	}

	req, err := http.NewRequest("POST", "/users/login", bytes.NewReader(marshalledSnap))

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

func getExistingUser(router *router.Router, username string) (int, UserProfile, error) {
	req, err := http.NewRequest("GET", "/users/"+username, &bytes.Reader{})

	if err != nil {
		return 0, UserProfile{}, err
	}

	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := UserProfile{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return 0, UserProfile{}, err
	}
	return recorder.Code, result, nil
}

func getRegisterOptions(router *router.Router) (int, RegisterOptions, error) {
	req, err := http.NewRequest("GET", "/users/register", &bytes.Reader{})

	if err != nil {
		return 0, RegisterOptions{}, err
	}

	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := RegisterOptions{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	if err != nil {
		return 0, RegisterOptions{}, err
	}
	return recorder.Code, result, nil
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

func AssertUserProfileIsUser(t *testing.T, user User, location string, interests []string, profile UserProfile) {
	assert.Equal(t, user.FirstName, profile.FirstName)
	assert.Equal(t, user.LastName, profile.LastName)
	assert.Equal(t, user.UserName, profile.UserName)
	assert.Equal(t, user.Mail, profile.Mail)
	assert.Equal(t, location, profile.Location)

	assert.Equal(t, len(profile.Interests), len(interests))
	for i, interest := range interests {
		assert.Equal(t, interest, profile.Interests[i])
	}
}
