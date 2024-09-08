package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"users-service/src/router"
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

func getExistingUser(router *router.Router, username string) (int, UserProfile, error) {
	req, err := http.NewRequest("GET", "/users/username"+username, &bytes.Reader{})

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

func CheckUserProfileIsUser(user User, location string, interests []string, profile UserProfile) bool {
	if user.FirstName != profile.FirstName {
		return false
	} else if user.LastName != profile.LastName {
		return false
	} else if user.UserName != profile.UserName {
		return false
	} else if user.Mail != profile.Mail {
		return false
	} else if location != profile.Location {
		return false
	} else if len(interests) != len(profile.Interests) {
		return false
	} else {
		for i, interest := range interests {
			if interest != profile.Interests[i] {
				return false
			}
		}
	}
	return true
}
