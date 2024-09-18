package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	email := "edwardelric@yahoo.com"
	locationId := 0
	interestsIds := []int{0, 1}
	personalInfo := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elri3c:)",
		Location:  locationId,
	}

	location, interests := getLocationAndInterestsNames(registerOptions, locationId, interestsIds)
	userProfile, err := CreateValidUser(router, email, personalInfo, interestsIds)

	assert.Equal(t, err, nil)

	AssertUserPrivateProfileIsUser(t, email, personalInfo, location, interests, userProfile)
}

func TestCreateUserWithInvalidPassword(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	personalInfo := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "EpaPapi",
		Location:  0,
	}

	code, response, err := CreateUserWithInvalidPersonalInfo(router, email, personalInfo)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "password")
	assertRegisterInstancePattern(t, "personal-info", response.Instance)
}

func TestCreateUserWithFirstAndLastNameTooLong(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	personalInfo := UserPersonalInfo{
		FirstName: "samdasdoasm,diop kiopsasdaskdopaskEdwardasdasdasdasdsdsadsasdasdasdasdasdasdasdasasdasdasasdasduashjudjashuasjdiasjdio",
		LastName:  "Lorem ipsum no se como sigue esto no se hablar latin juju asndjasndjasndjasndjkasnjsnadjkasndjkanjkasdffddfaafdafaasasd",
		UserName:  "EdwardoElric",
		Password:  "EpaPap#48efdwi",
		Location:  0,
	}

	code, response, err := CreateUserWithInvalidPersonalInfo(router, email, personalInfo)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 2)
	assertRegisterInstancePattern(t, "personal-info", response.Instance)
}

func TestCreateUserWithUsernameAndPasswordTooShort(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	personalInfo := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "As",
		Password:  "O^1i",
		Location:  0,
	}

	code, response, err := CreateUserWithInvalidPersonalInfo(router, email, personalInfo)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 2)
	assertRegisterInstancePattern(t, "personal-info", response.Instance)
}

func TestCreateUserWithInvalidLocation(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	code, registerOptions, err := getRegisterOptions(router)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)

	email := "capo@gmail.com"
	locationIndex := len(registerOptions.Locations) //invalid
	personalInfo := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$Elr3ic:)",
		Location:  locationIndex,
	}

	code, response, err := CreateUserWithInvalidPersonalInfo(router, email, personalInfo)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "location")
	assertRegisterInstancePattern(t, "personal-info", response.Instance)
}

func TestCreateUserWithNotExistingInterests(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	code, registerOptions, err := getRegisterOptions(router)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusOK)

	email := "capo@gmail.com"
	interstsIds := []int{len(registerOptions.Interests) - 1, len(registerOptions.Interests)}
	user := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$El3ric:)",
		Location:  0,
	}

	code, response, err := CreateUserWithInvalidInterests(router, email, user, interstsIds)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "interests")
	assertRegisterInstancePattern(t, "interests", response.Instance)
}

func TestCreateUserWithRepeatedInterests(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	interstsIds := []int{0, 0, 1}
	user := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Edward$El3ric:)",
		Location:  0,
	}

	code, response, err := CreateUserWithInvalidInterests(router, email, user, interstsIds)

	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "interests")
	assertRegisterInstancePattern(t, "interests", response.Instance)
}

func TestCreateUserWithInvalidMail(t *testing.T) {
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
	res := ValidationErrorResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	assert.Equal(t, err, nil)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
	assert.Equal(t, res.Title, "validation error")
	assert.Equal(t, len(res.Errors), 1)
	assert.Equal(t, res.Errors[0].Field, "email")
	assert.Equal(t, res.Instance, "/users/resolver")
}

func TestCreateUserWithMailThatExists(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	interestsIds := []int{0}
	personalInfo := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardElric",
		Password:  "Holaa&2dS",
		Location:  0,
	}

	_, err = CreateValidUser(router, email, personalInfo, interestsIds)
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
	res := ResolverResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	assert.Equal(t, err, nil)

	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(t, res.NextAuthStep, LoginAuthStep)
}

func TestCreateUserWithUsernameThatExistsWithDifferentCase(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "capo@gmail.com"
	interestsIds := []int{0}
	personalInfo := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Holaa&2dS",
		Location:  0,
	}

	_, err = CreateValidUser(router, email, personalInfo, interestsIds)
	assert.Equal(t, err, nil)

	email = "bestia@gmail.com"
	personalInfo = UserPersonalInfo{
		FirstName: "asda",
		LastName:  "Elrasdasdic",
		UserName:  "edWardoElrIc",
		Password:  "askdo02d(S",
		Location:  0,
	}

	code, response, err := CreateUserWithInvalidPersonalInfo(router, email, personalInfo)

	fmt.Println(response)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusBadRequest)
	assert.Equal(t, response.Title, "validation error")
	assert.Equal(t, len(response.Errors), 1)
	assert.Equal(t, response.Errors[0].Field, "username")
	assertRegisterInstancePattern(t, "personal-info", response.Instance)
}
