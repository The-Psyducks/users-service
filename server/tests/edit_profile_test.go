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

func setUpEditProfileTests() (testRouter *router.Router, user1 UserPrivateProfile, user1Password string, user2 UserPrivateProfile, user2Password string){
    var err error
    
    testRouter, err = router.CreateRouter()
    if err != nil {
        panic("Failed to create router: " + err.Error())
    }

    email := "edwardo@elric.com"
    locationId := 0
    interestsIds := []int{0, 1}
	user1Password = "Edward$El1ric:)"
    user := UserPersonalInfo{
        FirstName: "Edward",
        LastName:  "Elric",
        UserName:  "EdwardoElric",
        Password:  user1Password,
        Location:  locationId,
    }
    user1, err = CreateValidUser(testRouter, email, user, interestsIds)

    if err != nil {
        panic("Failed to create user1: " + err.Error())
    }

    email = "MonkeCAapoo@elric.com"
    locationId = 1
    interestsIds = []int{0, 1}
	user2Password = "Edward$El1ric:)"
    user = UserPersonalInfo{
        FirstName: "Monke",
        LastName:  "Unga",
        UserName:  "UngaUnga",
        Password:  user2Password,
        Location:  locationId,
    }
    user2, err = CreateValidUser(testRouter, email, user, interestsIds)
    if err != nil {
        panic("Failed to create user2: " + err.Error())
    }
    return testRouter, user1, user1Password, user2, user2Password
}

func TestEditProfile(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpEditProfileTests()
	options, err := getRegisterOptions(testRouter)
	assert.Equal(t, err, nil)

	updatedProfile := EditUserProfileRequest {
		FirstName: "Alphonse",
		LastName: "Elric",
		Username: "AlphonseElric",
		Location: len(options.Locations) - 1,
		Interests: []int{0, len(options.Interests) - 1},
	}
	logInRequest := LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}
	resp, err := LoginValidUser(testRouter, logInRequest)
	assert.Equal(t, err, nil)

	newUser, err := editValidUserProfile(testRouter, resp.AccessToken, updatedProfile)
	assert.Equal(t, err, nil)
	
	assert.Equal(t, newUser.FirstName, updatedProfile.FirstName)
	assert.Equal(t, newUser.LastName, updatedProfile.LastName)
	assert.Equal(t, newUser.UserName, updatedProfile.Username)
	assertInterestsNamesAreCorrectIds(t, options, updatedProfile.Interests, newUser.Interests)
	assertLocationNameIsCorrectId(t, options, updatedProfile.Location, newUser.Location)
}

func TestEditProfileChangesProfile(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpEditProfileTests()
	options, err := getRegisterOptions(testRouter)
	assert.Equal(t, err, nil)

	updatedProfile := EditUserProfileRequest {
		FirstName: "Alphonse",
		LastName: "Elric",
		Username: "AlphonseElric",
		Location: len(options.Locations) - 1,
		Interests: []int{0, len(options.Interests) - 1},
	}
	logInRequest := LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}
	resp, err := LoginValidUser(testRouter, logInRequest)
	assert.Equal(t, err, nil)

	newUser, err := editValidUserProfile(testRouter, resp.AccessToken, updatedProfile)
	assert.Equal(t, err, nil)

	userProfile, err := getOwnProfile(testRouter, newUser.Id.String(), resp.AccessToken)
	assert.Equal(t, err, nil)

	assertPrivateUsersAreEqual(t, newUser, userProfile)
}

func TestEditUnexistingProfile(t *testing.T) {
	testRouter, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	updatedProfile := EditUserProfileRequest {
		FirstName: "Alphonse",
		LastName: "Elric",
		Username: "AlphonseElric",
		Location: 0,
		Interests: []int{0},
	}
	userInfo, err := json.Marshal(updatedProfile)
	assert.Equal(t, err, nil)
	req, err := http.NewRequest("PUT", "/users/profile", bytes.NewReader(userInfo))
	assert.Equal(t, err, nil)
	req.Header.Add("Authorization", "Bearer "+ "123923")
	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	testRouter.Engine.ServeHTTP(recorder, req)
	result := ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.Equal(t, err, nil)

	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
	assert.Equal(t, result.Title, "Unauthorized")
}

func TestEditProfileWithInvalidData(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpEditProfileTests()

	updatedProfile := EditUserProfileRequest {
		FirstName: "",
		LastName: "E",
		Username: "as",
		Location: 5000,
		Interests: []int{0, 9000},
	}
	logInRequest := LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}
	resp, err := LoginValidUser(testRouter, logInRequest)
	assert.Equal(t, err, nil)

	updatedProfile.FirstName = ""
	code, response, err := editInvalidUserProfile(testRouter, resp.AccessToken, updatedProfile)
	
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 5)
}
