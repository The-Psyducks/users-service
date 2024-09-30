package tests

import (
	"fmt"
	"time"
    "encoding/json"
    "net/http"
    "net/http/httptest"
	"testing"
	"users-service/src/router"
	"github.com/go-playground/assert/v2"
)

func setUpSearchTests() (testRouter *router.Router, user1 UserPrivateProfile, user1Password string, user2 UserPrivateProfile, user2Password string){
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
        UserName:  "MonkeElrico",
        Password:  user2Password,
        Location:  locationId,
    }
    user2, err = CreateValidUser(testRouter, email, user, interestsIds)
    if err != nil {
        panic("Failed to create user2: " + err.Error())
    }
    return testRouter, user1, user1Password, user2, user2Password
}

func TestSearchForNotExistingUsersReturnsNothing(t *testing.T) {
	testRouter, user1, user1Password, _, _ := setUpSearchTests()

    LoginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password,
    }

    response, err := LoginValidUser(testRouter, LoginRequest)
    assert.Equal(t, err, nil)

    searchResult, err := searchUsers(testRouter, "^^", response.AccessToken, 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(searchResult), 0)
}

func TestSearchForOneUserReturnsJustIt(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpSearchTests()

    LoginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password,
    }

    response, err := LoginValidUser(testRouter, LoginRequest)
    assert.Equal(t, err, nil)

    searchResult, err := searchUsers(testRouter, "edward", response.AccessToken, 2)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(searchResult), 1)
    assert.Equal(t, searchResult[0].Profile.Id, user1.Id)
}

func TestSearchForUsersWithUsernameAndNameReturnsInOrder(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpSearchTests()

    LoginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password,
    }

    response, err := LoginValidUser(testRouter, LoginRequest)
    assert.Equal(t, err, nil)

    searchResult, err := searchUsers(testRouter, "elr", response.AccessToken, 2)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(searchResult), 2)
    assert.Equal(t, searchResult[0].Profile.Id, user2.Id)
    assert.Equal(t, searchResult[1].Profile.Id, user1.Id)
}

func TestSearchForUsersWithLowPaginationLimitReturnAll(t *testing.T) {
    testRouter, user1, user1Password, user2, _ := setUpSearchTests()

    LoginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password,
    }

    response, err := LoginValidUser(testRouter, LoginRequest)
    assert.Equal(t, err, nil)

    searchResult, err := searchUsers(testRouter, "elr", response.AccessToken, 1)
    assert.Equal(t, err, nil)
    assert.Equal(t, len(searchResult), 2)
    assert.Equal(t, searchResult[0].Profile.Id, user2.Id)
    assert.Equal(t, searchResult[1].Profile.Id, user1.Id)
}
func TestSearchForUsersWithWhitespacedTextReturnsError(t *testing.T) {
    testRouter, user1, user1Password, _, _ := setUpSearchTests()

    LoginRequest := LoginRequest{
        Email: user1.Email,
        Password: user1Password,
    }

    loginResponse, err := LoginValidUser(testRouter, LoginRequest)
    assert.Equal(t, err, nil)

    timestamp := time.Unix(time.Now().Unix()+1, 0).UTC().Format(time.RFC3339Nano)
    url := fmt.Sprintf("/users/search?text=%s&time=%s&skip=%d&limit=%d", " ", timestamp, 0, 2)
    req, err := http.NewRequest("GET", url, nil)
    assert.Equal(t, err, nil)

    req.Header.Add("Authorization", "Bearer "+ loginResponse.AccessToken)
    recorder := httptest.NewRecorder()
    testRouter.Engine.ServeHTTP(recorder, req)
    errorResponse := ErrorResponse{}
    err = json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
    assert.Equal(t, err, nil)

    assert.Equal(t, recorder.Code, http.StatusBadRequest)
    assert.Equal(t, errorResponse.Instance, "/users/search")
}