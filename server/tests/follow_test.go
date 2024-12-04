package tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"

	"users-service/src/router"
	"users-service/tests/models"
	"users-service/tests/utils"
)

func setUpFollowTests() (testRouter *router.Router,user1 models.UserPrivateProfile, user1Password string, user2 models.UserPrivateProfile, user2Password string){
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

func TestFollowUserAddsItToFollowers(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpFollowTests()
	loginRequest := models.LoginRequest{
		Email: user1.Email,
		Password: user1Password, 
	}
	resp, err := utils.LoginValidUser(testRouter, loginRequest)
	assert.Equal(t, err, nil)
    
	err = utils.FollowValidUser(testRouter, user2.Id.String(), resp.AccessToken)
	assert.Equal(t, err, nil)

    followers, err := utils.GetFollowers(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(followers), 1)
    assert.Equal(t, followers[0].Follows, false)
    assert.Equal(t, followers[0].Profile.UserName, user1.UserName)
}

func TestFollowUserTwiceReturnsProperError(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpFollowTests()
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    err = utils.FollowValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    code, response, err := utils.FollowInvalidUser(testRouter, user2.Id.String(), resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, code, http.StatusBadRequest)
    assert.Equal(t, response.Title, "The user already follows this user")
}

func TestFollowMyselfReturnsProperError(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpFollowTests()
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    code, response, err := utils.FollowInvalidUser(testRouter, user1.Id.String(), resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, code, http.StatusBadRequest)
    assert.Equal(t, response.Title, "Can't follow yourself")
}

func TestFollowNonExistentUserReturnsProperError(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpFollowTests()
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    code, response, err := utils.FollowInvalidUser(testRouter, "4a28ea57-854b-4ee4-af34-61c998d8493e", resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, code, http.StatusNotFound)
    assert.Equal(t, response.Title, "User not found")
}

func TestUnfollowUserDeletesItFromFollowers(t *testing.T) {
    testRouter, user1, user1Password, user2, user2Password := setUpFollowTests()
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    err = utils.FollowValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    
    err = utils.UnfollowValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    
    loginRequest = models.LoginRequest{
        Email: user2.Email,
        Password: user2Password, 
    }
    resp, err = utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    followers, err := utils.GetFollowers(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(followers), 0)
}

func TestUnfollowNonFollowingUserReturnsProperError(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpFollowTests()
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    code, response, err := utils.UnfollowInvalidUser(testRouter, user2.Id.String(), resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, code, http.StatusBadRequest)
    assert.Equal(t, response.Title, "The user is not following this user")
}

func TestUnfollowNonExistentUserReturnsProperError(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpFollowTests()
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    code, response, err := utils.UnfollowInvalidUser(testRouter, "4a28ea57-854b-4ee4-af34-61c998d8493e", resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, code, http.StatusNotFound)
    assert.Equal(t, response.Title, "User not found")
}

func TestGetFollowersWhenEmptyIsEmptyList(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpFollowTests()
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    followers, err := utils.GetFollowers(testRouter, user1.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(followers), 0)
}

func TestGetFollowersReturnsCorrectly(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpFollowTests()
    email := "asjid@elric.com"
    locationId := 0
    interestsIds := []int{0, 1}
    userPersonalInfo := models.UserPersonalInfo{
        FirstName: "Edward",
        LastName:  "Elric",
        UserName:  "Easokp",
        Password:  "sdaji@34fdasD",
        Location:  locationId,
    }
    user3, err := utils.CreateValidUser(testRouter, email, userPersonalInfo, interestsIds)
    assert.Equal(t, err, nil)

    loginRequest := models.LoginRequest{
        Email: user3.Email,
        Password: userPersonalInfo.Password,
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    err = utils.FollowValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    loginRequest = models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err = utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    err = utils.FollowValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    followers, err := utils.GetFollowers(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(followers), 2)
    assert.Equal(t, followers[0].Follows, false)
    assert.Equal(t, followers[1].Follows, false)
    profiles := []models.FollowUserProfile{{Follows: false, Profile: utils.PrivateUserToPublic(user1)}, {Follows: false, Profile: utils.PrivateUserToPublic(user3)}}
    utils.AssertListsAreEqual(t, followers, profiles)
}

func TestGetFollowingWhenEmptyReturnsEmptyList(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpFollowTests()
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    following, err := utils.GetFollowing(testRouter, user1.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(following), 0)
}

func TestGetFollowingReturnsCorrectly(t *testing.T) {
    testRouter, user1, _, user2, _ := setUpFollowTests()
    email := "asjid@elric.com"
    locationId := 0
    interestsIds := []int{0, 1}
    userPersonalInfo := models.UserPersonalInfo{
        FirstName: "Edward",
        LastName:  "Elric",
        UserName:  "Easokp",
        Password:  "sdaji@34fdasD",
        Location:  locationId,
    }
    user3, err := utils.CreateValidUser(testRouter, email, userPersonalInfo, interestsIds)
    assert.Equal(t, err, nil)
    
    loginRequest := models.LoginRequest{
        Email: user3.Email,
        Password: userPersonalInfo.Password,
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    err = utils.FollowValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    err = utils.FollowValidUser(testRouter, user1.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    following, err := utils.GetFollowing(testRouter, user3.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(following), 2)
    assert.Equal(t, following[0].Follows, true)
    assert.Equal(t, following[1].Follows, true)
    profiles := []models.FollowUserProfile{{Follows: true, Profile: utils.PrivateUserToPublic(user2)}, {Follows: true, Profile: utils.PrivateUserToPublic(user1)}}
    utils.AssertListsAreEqual(t, following, profiles)
}

func TestGetFollowersForNonFollowingUserReturnsProperError(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpFollowTests()
    
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password,
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    response, err := utils.GetFollowersForInvalidUser(testRouter, user2.Id.String(), resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, response.Title, "The user is not following this user")
}

func TestGetFollowerForMutualFollowingUsersReturnsCorrectly(t *testing.T) {
    testRouter, user1, user1Password, user2, user2Password := setUpFollowTests()
    email := "asjid@elric.com"
    locationId := 0
    interestsIds := []int{0, 1}
    userPersonalInfo := models.UserPersonalInfo{
        FirstName: "Edward",
        LastName:  "Elric",
        UserName:  "Easokp",
        Password:  "sdaji@34fdasD",
        Location:  locationId,
    }
    user3, err := utils.CreateValidUser(testRouter, email, userPersonalInfo, interestsIds)
    assert.Equal(t, err, nil)

    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    err = utils.FollowValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    err = utils.FollowValidUser(testRouter, user3.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    loginRequest = models.LoginRequest{
        Email: user2.Email,
        Password: user2Password, 
    }
    resp, err = utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    err = utils.FollowValidUser(testRouter, user1.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    err = utils.FollowValidUser(testRouter, user3.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    followers, err := utils.GetFollowers(testRouter, user3.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(followers), 2)
    for _, follower := range followers {
        if follower.Profile.UserName == user1.UserName {
            assert.Equal(t, follower.Follows, true)
        } else if follower.Profile.UserName == user2.UserName {
            assert.Equal(t, follower.Follows, false)
        } else {
            t.Errorf("Unexpected follower: %s", follower.Profile.UserName)
        }
    }
}

// tests:
// 1. Get Followers for invalid user
// 2. Get following for invalid user

func TestGetFollowersForNotExistingUserReturnsProperError(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpFollowTests()
    
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password,
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    response, err := utils.GetFollowersForInvalidUser(testRouter, "4a28ea57-854b-4ee4-af34-61c998d8493e", resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, response.Status, http.StatusNotFound)
}

func TestGetFollowingForNotExistingUserReturnsProperError(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpFollowTests()
    
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password,
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    response, err := utils.GetFollowingForInvalidUser(testRouter, "4a28ea57-854b-4ee4-af34-61c998d8493e", resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, response.Status, http.StatusNotFound)
}

func TestGetAmountOfFollowersInTimeRangeReturnsCorrecly(t *testing.T) {
    testRouter, user1, user1Password, user2, user2Password := setUpFollowTests()
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    err = utils.FollowValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    
    loginRequest = models.LoginRequest{
        Email: user2.Email,
        Password: user2Password, 
    }
    resp, err = utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    startTime := time.Now().Add(-time.Second).UTC().Format(time.RFC3339Nano)
    endTime := time.Now().UTC().Format(time.RFC3339Nano)
    amount, err := utils.GetAmountOfFollowersInTimeRange(testRouter, resp.AccessToken, startTime, endTime)
    assert.Equal(t, err, nil)
    assert.Equal(t, amount, 1)
}

func TestGetAmountOfFollowersInTimeRangeReturnsNoneIfThereAreNoNewFollowers(t *testing.T) {
    testRouter, user1, user1Password, user2, user2Password := setUpFollowTests()
    loginRequest := models.LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    err = utils.FollowValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    
    loginRequest = models.LoginRequest{
        Email: user2.Email,
        Password: user2Password, 
    }
    resp, err = utils.LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    startTime := time.Now().UTC().Format(time.RFC3339Nano)
    endTime := time.Now().Add(time.Second).UTC().Format(time.RFC3339Nano)
    amount, err := utils.GetAmountOfFollowersInTimeRange(testRouter, resp.AccessToken, startTime, endTime)
    assert.Equal(t, err, nil)
    assert.Equal(t, amount, 0)
}
