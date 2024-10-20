package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"users-service/src/app_errors"
	"users-service/src/database"
	"users-service/src/database/register_options"
	"users-service/src/model"

	"github.com/google/uuid"
)

func (u *User) GetUserProfileById(userSessionId uuid.UUID, userSessionIsAdmin bool, id uuid.UUID) (model.UserProfileResponse, error) {
	userRecord, err := u.userDb.GetUserById(id)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return model.UserProfileResponse{}, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return model.UserProfileResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}

	if userSessionId == id || userSessionIsAdmin {
		return u.getPrivateProfile(userRecord)
	}
	return u.getPublicProfile(userRecord, userSessionId)
}

func (u *User) getAmountOfFollowersAndFollowing(user model.UserRecord) (int, int, error) {
	followers, err := u.userDb.GetAmountOfFollowers(user.Id)
	if err != nil {
		return 0, 0, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting amount of followers: %w", err))
	}

	following, err := u.userDb.GetAmountOfFollowing(user.Id)
	if err != nil {
		return 0, 0, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error getting amount of following: %w", err))
	}

	return followers, following, nil
}

func (u *User) getPrivateProfile(user model.UserRecord) (model.UserProfileResponse, error) {
	privateProfile, err := u.createUserPrivateProfileFromUserRecord(user)
	if err != nil {
		return model.UserProfileResponse{}, err
	}

	slog.Info("user Private profile retrieved succesfully", slog.String("userId", user.Id.String()))
	return model.UserProfileResponse{
		OwnProfile: true,
		Follows:    false,
		Profile:    privateProfile,
	}, nil
}

func (u *User) getPublicProfile(user model.UserRecord, session_user_id uuid.UUID) (model.UserProfileResponse, error) {
	profile, err := u.generateUserPublicProfileFromUserRecord(user)
	if err != nil {
		return model.UserProfileResponse{}, err
	}

	follows, err := u.userDb.CheckIfUserFollows(session_user_id, user.Id)
	if err != nil {
		return model.UserProfileResponse{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error checking if user follows: %w", err))
	}

	slog.Info("user Public profile retrieved succesfully", slog.String("userId", user.Id.String()))
	return model.UserProfileResponse{
		OwnProfile: false,
		Follows:    follows,
		Profile:    profile,
	}, nil
}

func (u *User) validateUpdateUserPrivateProfile(data model.UpdateUserPrivateProfileRequest, userRecord model.UserRecord) ([]model.ValidationError, error) {
	totalValErrors := []model.ValidationError{}

	if !strings.EqualFold(data.UserName, userRecord.UserName) {
		if valErrs, err := u.userValidator.ValidateUpdateUsername(data.UserName); err != nil {
			return []model.ValidationError{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error validating username: %w", err))
		} else if len(valErrs) > 0 {
			totalValErrors = append(totalValErrors, valErrs...)
			fmt.Println("username val errs:", valErrs)
		}
	}

	updateProfileData := model.UpdateUserPrivateProfileData{
		PicturePath: data.PicturePath,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		Location:    data.Location,
		Interests:   data.Interests,
	}

	if valErrs, err := u.userValidator.ValidateUpdatePrivateProfileData(updateProfileData); err != nil {
		return []model.ValidationError{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error validating user personal info: %w", err))
	} else if len(valErrs) > 0 {
		totalValErrors = append(totalValErrors, valErrs...)
		fmt.Println("personal info val errs:", valErrs)
	}

	return totalValErrors, nil
}

func (u *User) ModifyUserProfile(userSessionId uuid.UUID, data model.UpdateUserPrivateProfileRequest) (model.UserPrivateProfile, error) {
	userRecord, err := u.userDb.GetUserById(userSessionId)
	if err != nil {
		if errors.Is(err, database.ErrKeyNotFound) {
			return model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusNotFound, UsernameNotFound, err)
		}
		return model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error retrieving user: %w", err))
	}
	
	valErrs, err := u.validateUpdateUserPrivateProfile(data, userRecord)
	if err != nil {
		return model.UserPrivateProfile{}, err
	} else if len(valErrs) > 0 {
		fmt.Println("all val errs:", valErrs)
		return model.UserPrivateProfile{}, app_errors.NewAppValidationError(valErrs)
	}

	location := register_options.GetLocationName(data.Location)
	interests := extractInterestNamesFromValidIds(data.Interests)
	updateData := model.UpdateUserPrivateProfile{
		UserName:    data.UserName,
		PicturePath: data.PicturePath,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		Location:    location,
		Interests:   interests,
	}

	updatedUser, err := u.userDb.ModifyUser(userSessionId, updateData)
	if err != nil {
		return model.UserPrivateProfile{}, app_errors.NewAppError(http.StatusInternalServerError, InternalServerError, fmt.Errorf("error updating user profile: %w", err))
	}
	privateProfile, err := u.createUserPrivateProfileFromUserRecord(updatedUser)
	if err != nil {
		return model.UserPrivateProfile{}, err
	}

	slog.Info("user profile updated succesfully", slog.String("userId", userSessionId.String()))
	return privateProfile, nil
}
