package tests

import (
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"

	"users-service/src/router"
	"users-service/tests/models"
	"users-service/tests/utils"
)

func TestLoginUserReturnsSession(t *testing.T) {
	router, err := router.CreateRouter()

	assert.Equal(t, err, nil)

	email := "edwardelric@yahoo.com"
	interestsIds := []int{0, 1}
	user := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elri3c:)",
		Location:  0,
	}

	_, err = utils.CreateValidUser(router, email, user, interestsIds)

	assert.Equal(t, err, nil)

	login := models.LoginRequest{
		Email:    email,
		Password: user.Password,
	}

	_, err = utils.LoginValidUser(router, login)

	assert.Equal(t, err, nil)
}

func TestLoginNotExistingUserReturnsNotFoundError(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	login := models.LoginRequest{
		Email:    "AtsumuMiya@GOAT.com",
		Password: "InarizakiGOAT",
	}

	code, resp, err := utils.LoginInvalidUser(router, login)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
	assert.Equal(t, resp.Title, "Incorrect username or password")
	assert.Equal(t, resp.Instance, "/users/login")
}

func TestLoginUserWithInvalidPasswordReturnsNotFoundError(t *testing.T) {
	router, err := router.CreateRouter()

	assert.Equal(t, err, nil)

	email := "edwardelric@yahoo.com"
	interestsIds := []int{0, 1}
	user := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elri3c:)",
		Location:  0,
	}

	_, err = utils.CreateValidUser(router, email, user, interestsIds)

	assert.Equal(t, err, nil)

	login := models.LoginRequest{
		Email:    email,
		Password: "Edward$Elri3c:",
	}

	code, _, err := utils.LoginInvalidUser(router, login)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
}

func TestLoginUserStillInRegistryReturnsNotFoundError(t *testing.T) {
	router, err := router.CreateRouter()

	assert.Equal(t, err, nil)

	email := "hola@gmail.com"
	interestsIds := []int{0, 1}
	personalInfo := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elri3c:)",
		Location:  0,
	}

	res, err := utils.GetUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)

	id := res.Metadata.RegistrationId

	err = utils.SendEmailVerificationAndVerificateIt(router, id)
	assert.Equal(t, err, nil)

	err = utils.PutValidUserPersonalInfo(router, id, personalInfo)
	assert.Equal(t, err, nil)

	err = utils.PutValidInterests(router, id, interestsIds)
	assert.Equal(t, err, nil)

	login := models.LoginRequest{
		Email:    email,
		Password: personalInfo.Password,
	}

	code, _, err := utils.LoginInvalidUser(router, login)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)

}
