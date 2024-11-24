package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"users-service/src/router"
	"users-service/tests/models"
	"users-service/tests/utils"

	"github.com/go-playground/assert/v2"
)

func setUpEditProfileTests() (testRouter *router.Router, user1 models.UserPrivateProfile, user1Password string, user2 models.UserPrivateProfile, user2Password string){
    var err error
    
    testRouter, err = router.CreateRouter()
    if err != nil {
        panic("Failed to create router: " + err.Error())
    }

    email := "edwardo@elric.com"
    locationId := 0
    interestsIds := []int{0, 1}
	user1Password = "Edward$El1ric:)"
    user := models.UserPersonalInfo{
        FirstName: "Edward",
        LastName:  "Elric",
        UserName:  "EdwardoElric",
        Password:  user1Password,
        Location:  locationId,
    }
    user1, err = utils.CreateValidUser(testRouter, email, user, interestsIds)

    if err != nil {
        panic("Failed to create user1: " + err.Error())
    }

    email = "MonkeCAapoo@elric.com"
    locationId = 1
    interestsIds = []int{0, 1}
	user2Password = "Edward$El1ric:)"
    user = models.UserPersonalInfo{
        FirstName: "Monke",
        LastName:  "Unga",
        UserName:  "UngaUnga",
        Password:  user2Password,
        Location:  locationId,
    }
    user2, err = utils.CreateValidUser(testRouter, email, user, interestsIds)
    if err != nil {
        panic("Failed to create user2: " + err.Error())
    }
    return testRouter, user1, user1Password, user2, user2Password
}

func TestEditProfileModifiesIt(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpEditProfileTests()
	options, err := utils.GetRegisterOptions(testRouter)
	assert.Equal(t, err, nil)

	updatedProfile := models.EditUserProfileRequest {
		FirstName: "Alphonse",
		LastName: "Elric",
		Username: "AlphonseElric",
		Location: len(options.Locations) - 1,
		Interests: []int{0, len(options.Interests) - 1},
	}
	logInRequest := models.LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}
	resp, err := utils.LoginValidUser(testRouter, logInRequest)
	assert.Equal(t, err, nil)

	newUser, err := utils.EditValidUserProfile(testRouter, resp.AccessToken, updatedProfile)
	assert.Equal(t, err, nil)
	
	assert.Equal(t, newUser.FirstName, updatedProfile.FirstName)
	assert.Equal(t, newUser.LastName, updatedProfile.LastName)
	assert.Equal(t, newUser.UserName, updatedProfile.Username)
	utils.AssertInterestsNamesAreCorrectIds(t, options, updatedProfile.Interests, newUser.Interests)
	utils.AssertLocationNameIsCorrectId(t, options, updatedProfile.Location, newUser.Location)
}

func TestEditProfileChangesPrivateProfile(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpEditProfileTests()
	options, err := utils.GetRegisterOptions(testRouter)
	assert.Equal(t, err, nil)

	updatedProfile := models.EditUserProfileRequest {
		FirstName: "Alphonse",
		LastName: "Elric",
		Username: "AlphonseElric",
		Location: len(options.Locations) - 1,
		Interests: []int{0, len(options.Interests) - 1},
	}
	logInRequest := models.LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}
	resp, err := utils.LoginValidUser(testRouter, logInRequest)
	assert.Equal(t, err, nil)

	newUser, err := utils.EditValidUserProfile(testRouter, resp.AccessToken, updatedProfile)
	assert.Equal(t, err, nil)

	userProfile, err := utils.GetOwnProfile(testRouter, newUser.Id.String(), resp.AccessToken)
	assert.Equal(t, err, nil)

	utils.AssertPrivateUsersAreEqual(t, newUser, userProfile)
}

func TestEditUnexistingProfileReturnsProperError(t *testing.T) {
	testRouter, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	updatedProfile := models.EditUserProfileRequest {
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
	result := models.ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.Equal(t, err, nil)

	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
	assert.Equal(t, result.Title, "Unauthorized")
}

func TestEditProfileWithInvalidDataReturnsProperValidationError(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpEditProfileTests()

	updatedProfile := models.EditUserProfileRequest {
		FirstName: "",
		LastName: "E",
		Username: "as",
		Location: 5000,
		Interests: []int{0, 9000},
	}
	logInRequest := models.LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}
	resp, err := utils.LoginValidUser(testRouter, logInRequest)
	assert.Equal(t, err, nil)

	updatedProfile.FirstName = ""
	code, response, err := utils.EditInvalidUserProfile(testRouter, resp.AccessToken, updatedProfile)
	
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 5)
}
