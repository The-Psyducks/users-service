package tests

import (
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"

	"users-service/src/router"
)

func setUpTest() (testRouter *router.Router,user1 UserPrivateProfile, user1Password string, user2 UserPrivateProfile, user2Password string){
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

func TestFollowUser(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpTest()
	loginRequest := LoginRequest{
		Email: user1.Email,
		Password: user1Password, 
	}
	resp, err := LoginValidUser(testRouter, loginRequest)
	assert.Equal(t, err, nil)
    
	err = followValidUser(testRouter, user2.Id.String(), resp.AccessToken)
	assert.Equal(t, err, nil)

    followers, err := getFollowers(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(followers), 1)
    assert.Equal(t, followers[0].Follows, false)
    assert.Equal(t, followers[0].Profile.UserName, user1.UserName)
}

func TestFollowUserTwice(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpTest()
    loginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    err = followValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    code, response, err := followInvalidUser(testRouter, user2.Id.String(), resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, code, http.StatusBadRequest)
    assert.Equal(t, response.Title, "The user already follows this user")
}

func TestFollowMyself(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpTest()
    loginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    code, response, err := followInvalidUser(testRouter, user1.Id.String(), resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, code, http.StatusBadRequest)
    assert.Equal(t, response.Title, "Can't follow yourself")
}

func TestFollowNonExistentUser(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpTest()
    loginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    code, response, err := followInvalidUser(testRouter, "4a28ea57-854b-4ee4-af34-61c998d8493e", resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, code, http.StatusNotFound)
    assert.Equal(t, response.Title, "User not found")
}

func TestUnfollowUser(t *testing.T) {
    testRouter, user1, user1Password, user2, user2Password := setUpTest()
    loginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    err = followValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    
    err = unfollowValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    
    loginRequest = LoginRequest{
        Email: user2.Email,
        Password: user2Password, 
    }
    resp, err = LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    followers, err := getFollowers(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(followers), 0)
}

func TestUnfollowNonFollowingUser(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpTest()
    loginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    code, response, err := unfollowInvalidUser(testRouter, user2.Id.String(), resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, code, http.StatusBadRequest)
    assert.Equal(t, response.Title, "The user is not following this user")
}

func TestUnfollowNonExistentUser(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpTest()
    loginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    code, response, err := unfollowInvalidUser(testRouter, "4a28ea57-854b-4ee4-af34-61c998d8493e", resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, code, http.StatusNotFound)
    assert.Equal(t, response.Title, "User not found")
}

func TestGetFollowersWhenEmpty(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpTest()
    loginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    followers, err := getFollowers(testRouter, user1.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(followers), 0)
}

func TestGetFollowers(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpTest()
    email := "asjid@elric.com"
    locationId := 0
    interestsIds := []int{0, 1}
    userPersonalInfo := UserPersonalInfo{
        FirstName: "Edward",
        LastName:  "Elric",
        UserName:  "Easokp",
        Password:  "sdaji@34fdasD",
        Location:  locationId,
    }
    user3, err := CreateValidUser(testRouter, email, userPersonalInfo, interestsIds)
    assert.Equal(t, err, nil)

    loginRequest := LoginRequest{
        Email: user3.Email,
        Password: userPersonalInfo.Password,
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    err = followValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    loginRequest = LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err = LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    err = followValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    followers, err := getFollowers(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(followers), 2)
    assert.Equal(t, followers[0].Follows, false)
    assert.Equal(t, followers[1].Follows, false)
    profiles := []FollowUserProfile{{Follows: false, Profile: privateUserToPublic(user1)}, {Follows: false, Profile: privateUserToPublic(user3)}}
    assertListsAreEqual(t, followers, profiles)
}

func TestGetFollowingWhenEmpty(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpTest()
    loginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)

    following, err := getFollowing(testRouter, user1.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(following), 0)
}

func TestGetFollowing(t *testing.T) {
    testRouter, user1, _, user2, _ := setUpTest()
    email := "asjid@elric.com"
    locationId := 0
    interestsIds := []int{0, 1}
    userPersonalInfo := UserPersonalInfo{
        FirstName: "Edward",
        LastName:  "Elric",
        UserName:  "Easokp",
        Password:  "sdaji@34fdasD",
        Location:  locationId,
    }
    user3, err := CreateValidUser(testRouter, email, userPersonalInfo, interestsIds)
    assert.Equal(t, err, nil)
    
    loginRequest := LoginRequest{
        Email: user3.Email,
        Password: userPersonalInfo.Password,
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    err = followValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    err = followValidUser(testRouter, user1.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    following, err := getFollowing(testRouter, user3.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(following), 2)
    assert.Equal(t, following[0].Follows, true)
    assert.Equal(t, following[1].Follows, true)
    profiles := []FollowUserProfile{{Follows: true, Profile: privateUserToPublic(user2)}, {Follows: true, Profile: privateUserToPublic(user1)}}
    assertListsAreEqual(t, following, profiles)
}

func TestCanNotGetFollowersForNonFollowingUser(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpTest()
    
    loginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password,
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    response, err := getFollowersForInvalidUser(testRouter, user2.Id.String(), resp.AccessToken)
    
    assert.Equal(t, err, nil)
    assert.Equal(t, response.Title, "The user is not following this user")
}

func TestGetFollowerForMutualFollowingUsers(t *testing.T) {
    testRouter, user1, user1Password, user2, user2Password := setUpTest()
    email := "asjid@elric.com"
    locationId := 0
    interestsIds := []int{0, 1}
    userPersonalInfo := UserPersonalInfo{
        FirstName: "Edward",
        LastName:  "Elric",
        UserName:  "Easokp",
        Password:  "sdaji@34fdasD",
        Location:  locationId,
    }
    user3, err := CreateValidUser(testRouter, email, userPersonalInfo, interestsIds)
    assert.Equal(t, err, nil)

    loginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password, 
    }
    resp, err := LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    err = followValidUser(testRouter, user2.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    err = followValidUser(testRouter, user3.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    loginRequest = LoginRequest{
        Email: user2.Email,
        Password: user2Password, 
    }
    resp, err = LoginValidUser(testRouter, loginRequest)
    assert.Equal(t, err, nil)
    
    err = followValidUser(testRouter, user1.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)
    err = followValidUser(testRouter, user3.Id.String(), resp.AccessToken)
    assert.Equal(t, err, nil)

    followers, err := getFollowers(testRouter, user3.Id.String(), resp.AccessToken)
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
