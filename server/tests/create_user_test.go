package tests

import (
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"

	"users-service/src/router"
)

func TestCreateUser(t *testing.T) {
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
		Password:  "Edward$Elri3c:)",
		Mail:      "edwardelric@yahoo.com",
		Location:  locationId,
		Interests: interestsIds,
	}

	code, userProfile, err := CreateValidUser(router, user)
	location, interests := getLocationAndInterestsNames(registerOptions, locationId, interestsIds)

	AssertUserProfileIsUser(t, user, location, interests, userProfile)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusCreated)
}

func TestCreateUserWithInvalidLocation(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	code, registerOptions, err := getRegisterOptions(router)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)

	locationIndex := len(registerOptions.Locations) //invalid
	user := User{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elr3ic:)",
		Mail:      "capo@gmail.com",
		Location:  int(locationIndex),
		Interests: []int{0},
	}

	code, response, err := CreateInvalidUser(router, user)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, response.Instance, "/users/register")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "location")
}

func TestCreateUserWithInvalidInterests(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	code, registerOptions, err := getRegisterOptions(router)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)

	interstsIds := []int{len(registerOptions.Interests) - 1, len(registerOptions.Interests)}
	user := User{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$El3ric:)",
		Mail:      "capo@gmail.com",
		Location:  0,
		Interests: interstsIds,
	}

	code, response, err := CreateInvalidUser(router, user)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, response.Instance, "/users/register")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "interests")
}

func TestCreateUserWithInvalidPassword(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	user := User{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "EpaPapi",
		Mail:      "capo@gmail.com",
		Location:  0,
		Interests: []int{0},
	}

	code, response, err := CreateInvalidUser(router, user)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, response.Instance, "/users/register")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "password")
}

func TestCreateUserWithUsernameThatExists(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	user := User{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Holaa&2dS",
		Mail:      "capo@gmail.com",
		Location:  0,
		Interests: []int{0},
	}

	code, _, err := CreateValidUser(router, user)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusCreated)

	user = User{
		FirstName: "asda",
		LastName:  "Elrasdasdic",
		UserName:  "EdwardoElric",
		Password:  "askdo02d(S",
		Mail:      "bestia@gmail.com",
		Location:  0,
		Interests: []int{0},
	}

	code, response, err := CreateInvalidUser(router, user)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusConflict)
	assert.Equal(t, response.Title, "Username or mail already exists")
	assert.Equal(t, response.Instance, "/users/register")
}

func TestCreateUserWithMailThatExists(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	code, _, err := getRegisterOptions(router)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)

	user := User{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardElric",
		Password:  "Holaa&2dS",
		Mail:      "capo@gmail.com",
		Location:  0,
		Interests: []int{0},
	}

	code, _, err = CreateValidUser(router, user)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusCreated)

	user = User{
		FirstName: "Alphonse",
		LastName:  "Elric",
		UserName:  "AlphonseElric",
		Password:  "Alph4on$eElric:)",
		Mail:      "capo@gmail.com",
		Location:  0,
		Interests: []int{0},
	}

	code, response, err := CreateInvalidUser(router, user)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusConflict)
	assert.Equal(t, response.Title, "Username or mail already exists")
	assert.Equal(t, response.Instance, "/users/register")
}
