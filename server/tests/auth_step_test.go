package tests

import (
	"testing"

	"github.com/go-playground/assert/v2"

	"users-service/src/router"
	"users-service/tests/constants"
	"users-service/tests/models"
	"users-service/tests/utils"
)

func TestUserAuthStepStartsWithEmailVerification(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "monke@gmail.com"

	res, err := utils.GetUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)
	assert.Equal(t, res.NextAuthStep, constants.SignUpAuthStep)
	assert.Equal(t, res.Metadata.OnboardingStep, constants.EmailVerificationStep)
}

func TestGetUserAuthStepWhenItsPersonalInfo(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "monke1@gmail.com"

	res, err := utils.GetUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)

	id := res.Metadata.RegistrationId

	err = utils.SendEmailVerificationAndVerificateIt(router, id)
	assert.Equal(t, err, nil)

	res, err = utils.GetUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)
	assert.Equal(t, res.NextAuthStep, constants.SignUpAuthStep)
	assert.Equal(t, res.Metadata.OnboardingStep, constants.PersonalInfoStep)
}

func TestGetUserAuthStepWhenItsInterests(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "monke2@gmail.com"

	res, err := utils.GetUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)

	id := res.Metadata.RegistrationId

	err = utils.SendEmailVerificationAndVerificateIt(router, id)
	assert.Equal(t, err, nil)

	personalInfo := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Holaa&2dS",
		Location:  0,
	}

	err = utils.PutValidUserPersonalInfo(router, id, personalInfo)
	assert.Equal(t, err, nil)

	res, err = utils.GetUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)
	assert.Equal(t, res.NextAuthStep, constants.SignUpAuthStep)
	assert.Equal(t, res.Metadata.OnboardingStep, constants.InterestsStep)
}

func TestGetUserAuthStepWhenItsComplete(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	email := "monke3@gmail.com"

	res, err := utils.GetUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)

	id := res.Metadata.RegistrationId

	err = utils.SendEmailVerificationAndVerificateIt(router, id)
	assert.Equal(t, err, nil)

	personalInfo := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Holaa&2dS",
		Location:  0,
	}

	err = utils.PutValidUserPersonalInfo(router, id, personalInfo)
	assert.Equal(t, err, nil)

	intetestsIds := []int{0}
	err = utils.PutValidInterests(router, id, intetestsIds)
	assert.Equal(t, err, nil)

	res, err = utils.GetUserRegistryForSignUp(router, email)
	assert.Equal(t, err, nil)
	assert.Equal(t, res.NextAuthStep, constants.SignUpAuthStep)
	assert.Equal(t, res.Metadata.OnboardingStep, constants.CompleteStep)
}

func AddPersonalInfoToNotExistingRegistry(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	personalInfo := models.UserPersonalInfo{
		FirstName: "Edward",
		LastName:  "Elric",
		UserName:  "EdwardoElric",
		Password:  "Holaa&2dS",
		Location:  0,
	}

	id := "0f246321-2921-41ad-8168-2b905b77a93c"
	err = utils.PutValidUserPersonalInfo(router, id, personalInfo)
	assert.Equal(t, err, nil)
}

func AddInterestsInfoToNotExistingRegistry(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	intetestsIds := []int{0}

	id := "0f246321-2921-41ad-8168-2b905b77a93c"
	err = utils.PutValidInterests(router, id, intetestsIds)
	assert.Equal(t, err, nil)
}
