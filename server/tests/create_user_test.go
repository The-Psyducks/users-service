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

	locationIndex := 0
	interestIndex1 := 0
	interestIndex2 := len(registerOptions.Interests) - 1
	user := User{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elric:)",
		Mail:      "edwardelric@yahoo.com",
		Location:  locationIndex,
		Interests: []int{interestIndex1, interestIndex2},
	}

	code, userProfile, err := CreateValidUser(router, user)
	location := registerOptions.Locations[locationIndex].Name
	interests := []string{registerOptions.Interests[interestIndex1].Name, registerOptions.Interests[interestIndex2].Name}
	equals := CheckUserProfileIsUser(user, location, interests, userProfile)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusCreated)
	assert.Equal(t, equals, true)
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
		Password:  "Edward$Elric:)",
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

	interestIndex1 := len(registerOptions.Interests) - 1 //valid
	interestIndex2 := len(registerOptions.Interests)     //invalid
	user := User{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elric:)",
		Mail:      "capo@gmail.com",
		Location:  0,
		Interests: []int{interestIndex1, interestIndex2},
	}

	code, response, err := CreateInvalidUser(router, user)
	
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, response.Instance, "/users/register")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "interests_ids")
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
	assert.Equal(t, response.Title, "username or email already exists")
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
		Password:  "Alphon$eElric:)",
		Mail:      "capo@gmail.com",
		Location:  0,
		Interests: []int{0},
	}

	code, response, err := CreateInvalidUser(router, user)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusConflict)
	assert.Equal(t, response.Title, "username or email already exists")
	assert.Equal(t, response.Instance, "/users/register")
}
