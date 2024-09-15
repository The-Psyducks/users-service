package tests

import (
	"net/http"
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
	code, result, err := getNotExistingUser(router, username)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
	assert.Equal(t, result.Title, "User not found")
}

func TestGetUserThatIsInRegistry(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "hola@gmail.com"
	interestsIds := []int{0, 1}
	personalInfo := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password: "Edward$Elri3c:)",
		Location:  0,
	}

	res, err := getUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)

	id := res.Metadata.RegistrationId

	err = sendEmailVerificationAndVerificateIt(router, id)
	assert.Equal(t, err, nil)

	_, err = putUserPersonalInfo(router, id, personalInfo)
	assert.Equal(t, err, nil)

	_, err = putInterests(router, id, interestsIds)
	assert.Equal(t, err, nil)

	code, result, err := getNotExistingUser(router, personalInfo.UserName)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
	assert.Equal(t, result.Title, "User not found")
}