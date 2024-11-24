package tests

import (
	"net/http"
	"testing"
	"users-service/src/router"
	"users-service/tests/models"
	"users-service/tests/utils"

	"github.com/go-playground/assert/v2"
)

func setUpBlockTests() (testRouter *router.Router, user1 models.UserPrivateProfile, user1Password string, user2 models.UserPrivateProfile, user2Password string){
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
        UserName:  "MonkeElrico",
        Password:  user2Password,
        Location:  locationId,
    }
    user2, err = utils.CreateValidUser(testRouter, email, user, interestsIds)
    if err != nil {
		panic("Failed to create user2: " + err.Error())
    }
	return
}

func TestBlockUserWorksCorrectly(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpBlockTests()

	adminToken, err := utils.LoginAdmin()
	assert.Equal(t, err, nil)

	err = utils.BlockUser(testRouter, user1.Id.String(), "You are blocked", adminToken)
	assert.Equal(t, err, nil)

	LoginRequest := models.LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}

	code, _, err := utils.LoginInvalidUser(testRouter, LoginRequest)
	assert.Equal(t, err, nil)

	assert.Equal(t, code, http.StatusForbidden)
}

func TestBlockUserWithoutReasonReturns400(t *testing.T) {
	testRouter, user1, _, _, _ := setUpBlockTests()

	adminToken, err := utils.LoginAdmin()
	assert.Equal(t, err, nil)

	code, _, err := utils.BlockInvalidUser(testRouter, user1.Id.String(), "", adminToken)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
}

func TestBlockUserWithNonAdminUserReturns403(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpBlockTests()

	loginReq := models.LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}

	user, err := utils.LoginValidUser(testRouter, loginReq)
	assert.Equal(t, err, nil)

	code, _, err := utils.BlockInvalidUser(testRouter, user1.Id.String(), "You are blocked", user.AccessToken)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusForbidden)
}

func TestUnblockUserWorksCorrectly(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpBlockTests()

	adminToken, err := utils.LoginAdmin()
	assert.Equal(t, err, nil)

	err = utils.BlockUser(testRouter, user1.Id.String(), "You are blocked", adminToken)
	assert.Equal(t, err, nil)

	err = utils.UnblockUser(testRouter, user1.Id.String(), adminToken)
	assert.Equal(t, err, nil)

	LoginRequest := models.LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}

	_,  err = utils.LoginValidUser(testRouter, LoginRequest)
	assert.Equal(t, err, nil)
}

func TestUnblockUserWithNonAdminUserReturns403(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpBlockTests()

	loginReq := models.LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}

	user, err := utils.LoginValidUser(testRouter, loginReq)
	assert.Equal(t, err, nil)

	code, _, err := utils.UnblockInvalidUser(testRouter, user1.Id.String(), user.AccessToken)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusForbidden)
}

func TestUnblockNotBlockedUserDoesNotFail(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpBlockTests()

	adminToken, err := utils.LoginAdmin()
	assert.Equal(t, err, nil)

	LoginRequest := models.LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}

	user, err := utils.LoginValidUser(testRouter, LoginRequest)
	assert.Equal(t, err, nil)

	err = utils.UnblockUser(testRouter, user.Profile.Id.String(), adminToken)
	assert.Equal(t, err, nil)
}

