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

	user := User{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elri3c:)",
		Mail:      "edwardelric@yahoo.com",
		Location:  0,
		Interests: []int{0, 1},
	}

	code, _, err := CreateValidUser(router, user)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusCreated)

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
	
	user := User{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elri3c:)",
		Mail:      "edwardelric@yahoo.com",
		Location:  0,
		Interests: []int{0, 1},
	}
	
	code, _, err := CreateValidUser(router, user)
	
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusCreated)
	
	login := LoginRequest{
		UserName: "AtsumuMiya",
		Password: "InarizakiGOAT",
	}
	
	code, resp, err := LoginInvalidUser(router, login)
	
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusUnauthorized)
	assert.Equal(t, resp.Title, "Incorrect username or password")
	assert.Equal(t, resp.Instance, "/users/login")
}