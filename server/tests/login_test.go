package tests

import (
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"

	"users-service/src/router"
)

func TestLoginUser(t *testing.T) {
	router, err := router.CreateRouter()

	assert.Equal(t, err, nil)

	email := "edwardelric@yahoo.com"
	interestsIds := []int{0, 1}
	user := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elri3c:)",
		Location:  0,
	}

	_, err = CreateValidUser(router, email, user, interestsIds)

	assert.Equal(t, err, nil)

	login := LoginRequest{
		UserName: user.UserName,
		Password: user.Password,
	}

	code, resp, err := LoginValidUser(router, login)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, resp.Valid, true)
}

func TestLoginNotExistingUser(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)
	
	login := LoginRequest{
		UserName: "AtsumuMiya",
		Password: "InarizakiGOAT",
	}
	
	code, resp, err := LoginInvalidUser(router, login)
	
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
	assert.Equal(t, resp.Title, "Incorrect username or password")
	assert.Equal(t, resp.Instance, "/users/login")
}

func TestLoginUserWithInvalidPassword(t *testing.T) {
	router, err := router.CreateRouter()

	assert.Equal(t, err, nil)

	email := "edwardelric@yahoo.com"
	interestsIds := []int{0, 1}
	user := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elri3c:)",
		Location:  0,
	}

	_, err = CreateValidUser(router, email, user, interestsIds)

	assert.Equal(t, err, nil)

	login := LoginRequest{
		UserName: user.UserName,
		Password: "Edward$Elri3c:",
	}

	code, _, err := LoginInvalidUser(router, login)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
}

func TestLoginUserStillInRegistry(t *testing.T) {
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

	login := LoginRequest{
		UserName: personalInfo.UserName,
		Password: personalInfo.Password,
	}

	code, _, err := LoginInvalidUser(router, login)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
	
}