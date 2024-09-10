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

	locationId := 0
	interestsIds := []int{0, 1}
	user := User{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$El1ric:)",
		Mail:      "edwardo@elric.com",
		Location:  locationId,
		Interests: interestsIds,
	}

	code, _, err = CreateValidUser(router, user)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusCreated)

	code, userProfile, err := getExistingUser(router, user.UserName)
	location, interests := getLocationAndInterestsNames(registerOptions, locationId, interestsIds)
	AssertUserProfileIsUser(t, user, location, interests, userProfile)

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
