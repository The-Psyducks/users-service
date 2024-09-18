package tests

import (
	"bytes"
	"testing"
	"net/http"
	"encoding/json"
	"net/http/httptest"
	"github.com/go-playground/assert/v2"

	"users-service/src/router"
)

func TestGetOwnUserProfile(t *testing.T) {
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

	code, response, err := createAndLoginUser(router, email, user, interestsIds)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)

	code, userProfile, err := getExistingUser(router, user.UserName, response.AccessToken)
	location, interests := getLocationAndInterestsNames(registerOptions, locationId, interestsIds)
	AssertUserPrivateProfileIsUser(t, email, user, location, interests, userProfile)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)
}

//get withoyut token
func TestGetUserProfileWithoutToken(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

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

	req, err := http.NewRequest("GET", "/users/"+user.UserName, &bytes.Reader{})
	assert.Equal(t, err, nil)

	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.Equal(t, err, nil)

	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
	assert.Equal(t, result.Title, "Authorization header is required")
}


func TestGetNotExistingUserProfile(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "edwardos@elric.com"
	interestsIds := []int{0, 1}
	user := UserPersonalInfo{
		FirstName: "Edwarsd",
		LastName:  "Elrsic",
		UserName:  "EdwasrdoElric",
		Password:  "Edward$El1ric:)",
		Location:  0,
	}

	code, response, err := createAndLoginUser(router, email, user, interestsIds)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)

	username := "monkeCrack"
	code, result, err := getNotExistingUser(router, username, response.AccessToken)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
	assert.Equal(t, result.Title, "User not found")
}

func TestGetUserThatIsInRegistry(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "hola@gmail.com"
	interestsIds := []int{0, 1}
	user := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elri3c:)",
		Location:  0,
	}

	code, response, err := createAndLoginUser(router, email, user, interestsIds)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)

	email = "holasa@gmail.com"
	interestsIds = []int{0, 1}
	user = UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoasElric",
		Password:  "Edward$Elri3c:)",
		Location:  0,
	}

	res, err := getUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)

	id := res.Metadata.RegistrationId

	err = sendEmailVerificationAndVerificateIt(router, id)
	assert.Equal(t, err, nil)

	_, err = putUserPersonalInfo(router, id, user)
	assert.Equal(t, err, nil)

	_, err = putInterests(router, id, interestsIds)
	assert.Equal(t, err, nil)

	code, result, err := getNotExistingUser(router, user.UserName, response.AccessToken)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
	assert.Equal(t, result.Title, "User not found")
}
