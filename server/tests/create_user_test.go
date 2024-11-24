package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"

	"users-service/src/router"
	"users-service/tests/constants"
	"users-service/tests/models"
	"users-service/tests/utils"
)

func TestCreateUserCreatesAValidOne(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	registerOptions, err := utils.GetRegisterOptions(router)
	assert.Equal(t, err, nil)

	email := "edwardelric@yahoo.com"
	locationId := 0
	interestsIds := []int{0, 1}
	personalInfo := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elri3c:)",
		Location:  locationId,
	}

	location, interests := utils.GetLocationAndInterestsNames(registerOptions, locationId, interestsIds)
	userProfile, err := utils.CreateValidUser(router, email, personalInfo, interestsIds)

	assert.Equal(t, err, nil)

	utils.AssertUserPrivateProfileIsUser(t, email, personalInfo, location, interests, userProfile)
}

func TestCreateUserWithInvalidPasswordReturnsProperValidationError(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	personalInfo := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "EpaPapi",
		Location:  0,
	}

	code, response, err := utils.CreateUserWithInvalidPersonalInfo(router, email, personalInfo)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "password")
	utils.AssertRegisterInstancePattern(t, "personal-info", response.Instance)
}

func TestCreateUserWithFirstAndLastNameTooLongReturnsProperValidationError(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	personalInfo := models.UserPersonalInfo{
		FirstName: "samdasdoasm,diop kiopsasdaskdopaskEdwardasdasdasdasdsdsadsasdasdasdasdasdasdasdasasdasdasasdasduashjudjashuasjdiasjdio",
		LastName:  "Lorem ipsum no se como sigue esto no se hablar latin juju asndjasndjasndjasndjkasnjsnadjkasndjkanjkasdffddfaafdafaasasd",
		UserName:  "EdwardoElric",
		Password:  "EpaPap#48efdwi",
		Location:  0,
	}

	code, response, err := utils.CreateUserWithInvalidPersonalInfo(router, email, personalInfo)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 2)
	utils.AssertRegisterInstancePattern(t, "personal-info", response.Instance)
}

func TestCreateUserWithUsernameAndPasswordTooShortReturnsProperValidationErrors(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	personalInfo := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "As",
		Password:  "O^1i",
		Location:  0,
	}

	code, response, err := utils.CreateUserWithInvalidPersonalInfo(router, email, personalInfo)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 2)
	utils.AssertRegisterInstancePattern(t, "personal-info", response.Instance)
}

func TestCreateUserWithInvalidLocationReturnsProperValidationError(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	registerOptions, err := utils.GetRegisterOptions(router)
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	locationIndex := len(registerOptions.Locations) //invalid
	personalInfo := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elr3ic:)",
		Location:  locationIndex,
	}

	code, response, err := utils.CreateUserWithInvalidPersonalInfo(router, email, personalInfo)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "location")
	utils.AssertRegisterInstancePattern(t, "personal-info", response.Instance)
}

func TestCreateUserWithNotExistingInterestsReturnsProperValidationError(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	registerOptions, err := utils.GetRegisterOptions(router)
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	interstsIds := []int{len(registerOptions.Interests) - 1, len(registerOptions.Interests)}
	user := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$El3ric:)",
		Location:  0,
	}

	code, response, err := utils.CreateUserWithInvalidInterests(router, email, user, interstsIds)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "interests")
	utils.AssertRegisterInstancePattern(t, "interests", response.Instance)
}

func TestCreateUserWithRepeatedInterestsReturnsProperValidationError(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	interstsIds := []int{0, 0, 1}
	user := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$El3ric:)",
		Location:  0,
	}

	code, response, err := utils.CreateUserWithInvalidInterests(router, email, user, interstsIds)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "interests")
	utils.AssertRegisterInstancePattern(t, "interests", response.Instance)
}

func TestCreateUserWithInvalidMailReturnsProperValidationError(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	payload := map[string]string{
		"email": "como andamos king",
	}
	marshalledInfo, err := json.Marshal(payload)
	assert.Equal(t, err, nil)

	req, err := http.NewRequest("POST", "/users/resolver", bytes.NewReader(marshalledInfo))
	assert.Equal(t, err, nil)

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	res := models.ValidationErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	assert.Equal(t, err, nil)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
	assert.Equal(t, res.Title, "validation error")
	assert.Equal(t, len(res.Errors), 1)
	assert.Equal(t, res.Errors[0].Field, "email")
	assert.Equal(t, res.Instance, "/users/resolver")
}

func TestResolveUserWithMailThatExistsReturnsLogInStep(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	interestsIds := []int{0}
	personalInfo := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardElric",
		Password:  "Holaa&2dS",
		Location:  0,
	}

	_, err = utils.CreateValidUser(router, email, personalInfo, interestsIds)
	assert.Equal(t, err, nil)

	payload := map[string]string{
		"email": email,
	}
	marshalledInfo, err := json.Marshal(payload)
	assert.Equal(t, err, nil)

	req, err := http.NewRequest("POST", "/users/resolver", bytes.NewReader(marshalledInfo))
	assert.Equal(t, err, nil)

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)
	res := models.ResolverResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	assert.Equal(t, err, nil)

	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(t, res.NextAuthStep, constants.LoginAuthStep)
}

func TestCreateUserWithUsernameThatExistsWithDifferentCaseReturnsProperValidationError(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	interestsIds := []int{0}
	personalInfo := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Holaa&2dS",
		Location:  0,
	}

	_, err = utils.CreateValidUser(router, email, personalInfo, interestsIds)
	assert.Equal(t, err, nil)

	email = "bestia@gmail.com"
	personalInfo = models.UserPersonalInfo{
		FirstName: "asda",
		LastName:  "Elrasdasdic",
		UserName:  "edWardoElrIc",
		Password:  "askdo02d(S",
		Location:  0,
	}

	code, response, err := utils.CreateUserWithInvalidPersonalInfo(router, email, personalInfo)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "username")
	utils.AssertRegisterInstancePattern(t, "personal-info", response.Instance)
}
