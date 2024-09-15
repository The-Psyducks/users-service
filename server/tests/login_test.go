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

func TestLoginInvalidUser(t *testing.T) {
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
		UserName: "AtsumuMiya",
		Password: "InarizakiGOAT",
	}

	code, resp, err := LoginInvalidUser(router, login)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
	assert.Equal(t, resp.Title, "Incorrect username or password")
	assert.Equal(t, resp.Instance, "/users/login")
}
