package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"

	"users-service/src/router"
)

func TestGetUser(t *testing.T) {
	router, err := router.CreateRouter()

	assert.Equal(t, err, nil)

	code, registerOptions, err := getRegisterOptions(router)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)

	email := "edwardo@elric.com"
	locationId := 0
	interestsIds := []int{0, 1}
	user := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$El1ric:)",
		Location:  locationId,
	}

	_, err = CreateValidUser(router, email, user, interestsIds)

	assert.Equal(t, err, nil)

	code, userProfile, err := getExistingUser(router, user.UserName)
	location, interests := getLocationAndInterestsNames(registerOptions, locationId, interestsIds)
	AssertUserProfileIsUser(t, email, user, location, interests, userProfile)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)
}

func TestGetNotExistingUser(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	username := "monkeCrack"
	req, err := http.NewRequest("GET", "/users/username"+username, &bytes.Reader{})

	assert.Equal(t, err, nil)

	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)

	assert.Equal(t, err, nil)
	assert.Equal(t, recorder.Code, http.StatusNotFound)
	assert.Equal(t, result.Title, "User not found")
}
