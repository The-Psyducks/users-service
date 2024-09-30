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

func TestGetOwnUserProfileReturnsPrivateProfile(t *testing.T) {
	router, err := router.CreateRouter()

	assert.Equal(t, err, nil)

	registerOptions, err := getRegisterOptions(router)
	assert.Equal(t, err, nil)

	email := "edwardo@elric.com"
	locationId := 0
	interestsIds := []int{0, 1}
	user := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$El1ric:)",
		Location:  locationId,
	}

	response, err := createAndLoginUser(router, email, user, interestsIds)
	assert.Equal(t, err, nil)

	userProfile, err := getOwnProfile(router, response.Profile.Id.String(), response.AccessToken)
	location, interests := getLocationAndInterestsNames(registerOptions, locationId, interestsIds)
	AssertUserPrivateProfileIsUser(t, email, user, location, interests, userProfile)

	assert.Equal(t, err, nil)
}

func TestGetAnotherUserProfileReturnsPublicProfile(t *testing.T) {
	router, err := router.CreateRouter()

	assert.Equal(t, err, nil)

	registerOptions, err := getRegisterOptions(router)
	assert.Equal(t, err, nil)

	email := "edwardo@elric.com"
	locationId := 0
	interestsIds := []int{0, 1}
	user := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$El1ric:)",
		Location:  locationId,
	}
	_, err = CreateValidUser(router, email, user, interestsIds)
	assert.Equal(t, err, nil)

	email = "MonkeCAapoo@elric.com"
	locationId = 1
	interestsIds = []int{0, 1}
	user = UserPersonalInfo{
		FirstName: "Monke",
		LastName:  "Unga",
		UserName:  "UngaUnga",
		Password:  "Edward$Esl1ric:)",
		Location:  locationId,
	}

	response2, err := createAndLoginUser(router, email, user, interestsIds)
	assert.Equal(t, err, nil)

	userProfile, err := getAnotherUserProfile(router, response2.Profile.Id.String(), response2.AccessToken)
	location, _ := getLocationAndInterestsNames(registerOptions, locationId, interestsIds)
	AssertUserPublicProfileIsUser(t, user, location, userProfile)

	assert.Equal(t, err, nil)
}

func TestGetUserProfileWithoutTokenReturnsAuthError(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "edwardo@elric.com"
	locationId := 0
	interestsIds := []int{0, 1}
	user := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$El1ric:)",
		Location:  locationId,
	}

	_, err = CreateValidUser(router, email, user, interestsIds)
	assert.Equal(t, err, nil)

	req, err := http.NewRequest("GET", "/users/"+user.UserName, &bytes.Reader{})
	assert.Equal(t, err, nil)

	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	result := ErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.Equal(t, err, nil)

	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
	assert.Equal(t, result.Title, "Unauthorized")
}

func TestGetNotExistingUserProfileReturnsProperError(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "edwardos@elric.com"
	interestsIds := []int{0, 1}
	user := UserPersonalInfo{
		FirstName: "Edwarsd",
		LastName:  "Elrsic",
		UserName:  "EdwasrdoElric",
		Password:  "Edward$El1ric:)",
		Location:  0,
	}

	response, err := createAndLoginUser(router, email, user, interestsIds)
	assert.Equal(t, err, nil)

	id := "4a28ea57-854b-4ee4-af34-61c998d8493e"
	result, err := getNotExistingUser(router, id, response.AccessToken)

	assert.Equal(t, err, nil)
	assert.Equal(t, result.Title, "User not found")
}

func TestGetUserThatIsInRegistryReturnsNotFoundError(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "hola@gmail.com"
	interestsIds := []int{0, 1}
	user := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elri3c:)",
		Location:  0,
	}

	response, err := createAndLoginUser(router, email, user, interestsIds)
	assert.Equal(t, err, nil)

	email = "holasa@gmail.com"
	interestsIds = []int{0, 1}
	user = UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoasElric",
		Password:  "Edward$Elri3c:)",
		Location:  0,
	}

	res, err := getUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)

	id := res.Metadata.RegistrationId

	err = sendEmailVerificationAndVerificateIt(router, id)
	assert.Equal(t, err, nil)

	err = putValidUserPersonalInfo(router, id, user)
	assert.Equal(t, err, nil)

	err = putValidInterests(router, id, interestsIds)
	assert.Equal(t, err, nil)

	result, err := getNotExistingUser(router, res.Metadata.RegistrationId, response.AccessToken)

	assert.Equal(t, err, nil)
	assert.Equal(t, result.Title, "User not found")
}
