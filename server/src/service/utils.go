package service

import (
	"fmt"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/constants"
	"users-service/src/database/register_options"
	"users-service/src/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func generateUserPersonalInfoRecordFromRequest(request model.UserPersonalInfoRequest) (*model.UserPersonalInfoRecord, error) {
	password, err := hashPassword(request.Password)
	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, "Internal server error", fmt.Errorf("error hashing password: %w", err))
	}

	return &model.UserPersonalInfoRecord{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		UserName:  request.UserName,
		Password:  password,
		Location:  register_options.GetLocationName(request.LocationId),
	}, nil
}

func generateUserRecordFromRegistryEntry(registry model.RegistryEntry) model.UserRecord {
	return model.UserRecord{
		UserName:  registry.PersonalInfo.UserName,
		FirstName: registry.PersonalInfo.FirstName,
		LastName:  registry.PersonalInfo.LastName,
		Email:     registry.Email,
		Password:  registry.PersonalInfo.Password,
		Location:  registry.PersonalInfo.Location,
		Interests: registry.Interests,
	}
}

func (u *User) createUserPrivateProfileFromUserRecord(record model.UserRecord) (model.UserPrivateProfile, error) {
	followers, following, err := u.getAmountOfFollowersAndFollowing(record)
	if err != nil {
		return model.UserPrivateProfile{}, err
	}

	return model.UserPrivateProfile{
		Id:        record.Id,
		UserName:  record.UserName,
		FirstName: record.FirstName,
		LastName:  record.LastName,
		Email:     record.Email,
		Location:  record.Location,
		Interests: record.Interests,
		Followers: followers,
		Following: following,
	}, nil
}

func (u *User) generateUserPublicProfileFromUserRecord(user model.UserRecord) (model.UserPublicProfile, error) {
	followers, following, err := u.getAmountOfFollowersAndFollowing(user)
	if err != nil {
		return model.UserPublicProfile{}, err
	}

	return model.UserPublicProfile{
		Id:        user.Id,
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Location:  user.Location,
		Followers: followers,
		Following: following,
	}, nil
}

func (u *User) getFollowStatusPublicProfilesFromUserRecords(userRecords []model.UserRecord, sessionUserId uuid.UUID) ([]model.UserPublicProfileWithFollowStatus, error) {
	profiles := make([]model.UserPublicProfileWithFollowStatus, 0, len(userRecords))
	for _, user := range userRecords {
		profile, err := u.generateUserPublicProfileFromUserRecord(user)
		if err != nil {
			return nil, err
		}
		follows, err := u.userDb.CheckIfUserFollows(sessionUserId, user.Id)
		if err != nil {
			return nil, app_errors.NewAppError(http.StatusInternalServerError, "Internal server error", fmt.Errorf("error checking if user follows: %w", err))
		}
		followProfile := model.UserPublicProfileWithFollowStatus{
			Follows: follows,
			Profile: profile,
		}
		profiles = append(profiles, followProfile)
	}
	return profiles, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", app_errors.NewAppError(http.StatusInternalServerError, "Internal server error", fmt.Errorf("error hashing password: %w", err))
	}

	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func extractInterestNamesFromValidIds(interests []int) []string {
	interestsNames := make([]string, len(interests))
	for i, interest := range interests {
		if name := register_options.GetInterestName(interest); name != "" {
			interestsNames[i] = name
		}
	}
	return interestsNames
}

func getStepForRegistryEntry(entry model.RegistryEntry) string {
	if !entry.EmailVerified {
		return constants.EmailVerificationStep
	}

	if entry.PersonalInfo.FirstName == "" {
		return constants.PersonalInfoStep
	}

	if len(entry.Interests) == 0 {
		return constants.InterestsStep
	}

	return constants.CompleteStep
}
