package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"users-service/src/router"

	"github.com/go-playground/assert/v2"
)

func setUpRecommendationTests() (testRouter *router.Router, user1 UserPrivateProfile, user1Password string, user2 UserPrivateProfile, user2Password string, user3 UserPrivateProfile, user3Password string, user4 UserPrivateProfile, user4Password string){
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

	email = "Slickback@gmail.com"
    locationId = 2
    interestsIds = []int{2, 3}
	user3Password = "Slickback1!"
    user = UserPersonalInfo{
        FirstName: "Joaquin",
        LastName:  "Pandolfi",
        UserName:  "Slickback",
        Password:  user3Password,
        Location:  locationId,
    }
    user3, err = CreateValidUser(testRouter, email, user, interestsIds)
    if err != nil {
        panic("Failed to create user3: " + err.Error())
    }

	email = "jinglebell@gmail.com"
    locationId = 0
    interestsIds = []int{0, 3}
	user3Password = "JingleBell1!"
    user = UserPersonalInfo{
        FirstName: "Martina",
        LastName:  "Lozano",
        UserName:  "JingleBell",
        Password:  user3Password,
        Location:  locationId,
    }
    user3, err = CreateValidUser(testRouter, email, user, interestsIds)
    if err != nil {
        panic("Failed to create user4: " + err.Error())
    }
    return testRouter, user1, user1Password, user2, user2Password, user3, user3Password, user4, user4Password
}

func TestGetRecommendationsWithInvalidUUIDReturnsUnauthorizedError(t *testing.T) {
    testRouter, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	timestamp := time.Unix(time.Now().Unix()+1, 0).UTC().Format(time.RFC3339Nano)
	url := fmt.Sprintf("/users/recommendations?&time=%s&skip=%d&limit=%d", timestamp, 0, 10)
	req, err := http.NewRequest("GET", url, nil)
	assert.Equal(t, err, nil)

	randomUUID := "e3901fcc-9287-4775-a484-f363f69cefd3"
	req.Header.Add("Authorization", "Bearer "+ randomUUID)
	recorder := httptest.NewRecorder()
	testRouter.Engine.ServeHTTP(recorder, req)
	
	newResult := ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &newResult)
	assert.Equal(t, err, nil)

	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
	assert.Equal(t, newResult.Status, http.StatusUnauthorized)
}

func TestGetRecommendationsForUserReturnsCorrectRecommendations(t *testing.T) {
	testRouter, user1, user1Password, user2,_,_,_,user4,_ := setUpRecommendationTests()

	LoginRequest := LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}

	response, err := LoginValidUser(testRouter, LoginRequest)
	assert.Equal(t, err, nil)

	recommendations, err := getAllUserRecommendations(testRouter, response.AccessToken, 10)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(recommendations), 2)
	assert.Equal(t, recommendations[0].Profile.Id, user4.Id)
	assert.Equal(t, recommendations[1].Profile.Id, user2.Id)
}

func TestGetRecommendationsForUserReturnsPaginatedRecommendations(t *testing.T) {
	testRouter, user1, user1Password, user2,_,_,_,user4,_ := setUpRecommendationTests()

	LoginRequest := LoginRequest{
		Email: user1.Email,
		Password: user1Password,
	}

	response, err := LoginValidUser(testRouter, LoginRequest)
	assert.Equal(t, err, nil)

	recommendations, err := getAllUserRecommendations(testRouter, response.AccessToken, 1)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(recommendations), 2)
	assert.Equal(t, recommendations[0].Profile.Id, user4.Id)
	assert.Equal(t, recommendations[1].Profile.Id, user2.Id)
}


