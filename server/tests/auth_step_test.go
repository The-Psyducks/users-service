package tests

import (
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"

	"users-service/src/router"
)

func TestUserAuthStepStartsWithEmailVerifiaction(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "monke@gmail.com"

	res, err := getUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)
	assert.Equal(t, res.NextAuthStep, SignUpAuthStep)
	assert.Equal(t, res.Metadata.OnboardingStep, EmailVerificationStep)
}

func TestGetUserAuthStepWhenItsPersonalInfo(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "monke@gmail.com"

	res, err := getUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)

	id := res.Metadata.RegistrationId

	err = sendEmailVerificationAndVerificateIt(router, id)
	assert.Equal(t, err, nil)

	res, err = getUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)
	assert.Equal(t, res.NextAuthStep, SignUpAuthStep)
	assert.Equal(t, res.Metadata.OnboardingStep, PersonalInfoStep)
}

func TestGetUserAuthStepWhenItsInterests(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "monke@gmail.com"

	res, err := getUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)

	id := res.Metadata.RegistrationId

	err = sendEmailVerificationAndVerificateIt(router, id)
	assert.Equal(t, err, nil)

	personalInfo := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Holaa&2dS",
		Location:  0,
	}

	code, err := putUserPersonalInfo(router, id, personalInfo)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNoContent)

	res, err = getUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)
	assert.Equal(t, res.NextAuthStep, SignUpAuthStep)
	assert.Equal(t, res.Metadata.OnboardingStep, InterestsStep)
}

func TestGetUserAuthStepWhenItsComplete (t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "monke@gmail.com"

	res, err := getUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)

	id := res.Metadata.RegistrationId

	err = sendEmailVerificationAndVerificateIt(router, id)
	assert.Equal(t, err, nil)

	personalInfo := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Holaa&2dS",
		Location:  0,
	}

	code, err := putUserPersonalInfo(router, id, personalInfo)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNoContent)

	intetestsIds := []int{0}
	code, err = putInterests(router, id, intetestsIds)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNoContent)

	res, err = getUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)
	assert.Equal(t, res.NextAuthStep, SignUpAuthStep)
	assert.Equal(t, res.Metadata.OnboardingStep, CompleteStep)
}

func AddPersonalInfoToNotExistingRegistry(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	personalInfo := UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Holaa&2dS",
		Location:  0,
	}

	id := "0f246321-2921-41ad-8168-2b905b77a93c"
	code, err := putUserPersonalInfo(router, id, personalInfo)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
}

func AddInterestsInfoToNotExistingRegistry(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	intetestsIds := []int{0}

	id := "0f246321-2921-41ad-8168-2b905b77a93c"
	code, err := putInterests(router, id, intetestsIds)
	assert.Equal(t, err, nil)
	assert.Equal(t, code, http.StatusNotFound)
}