package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ResolveRequest struct {
	Email        string   `json:"email" validate:"required"`
	ProviderData Provider `json:"provider"`
}

type Provider struct {
	Name     string      `json:"type" validate:"required"`
	Metadata json.RawMessage `json:"metadata" validate:"required"`
}

type GoogleAuthMetadata struct {
	FirebaseTokenId	string	`json:"firebase_token_id" validate:"required"`
}

type ResolveResponse struct {
	NextAuthStep string      `json:"next_auth_step" validate:"required"`
	Metadata     interface{} `json:"metadata" validate:"required"`
}

type UserPersonalInfoRequest struct {
    FirstName  string `json:"first_name" validate:"required,firstnamevalidator"`
    LastName   string `json:"last_name" validate:"required,lastnamevalidator"`
    UserName   string `json:"username" validate:"required,usernamevalidator"`
    Password   string `json:"password" validate:"required,passwordvalidator"`
    LocationId int    `json:"location" validate:"locationvalidator"`
}

type UserPersonalInfoRecord struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	UserName  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Location  string `json:"location" validate:"required"`
}

type RegistryEntry struct {
	Id            uuid.UUID              `json:"id" validate:"required"`
	Email         string                 `json:"email" validate:"required"`
	EmailVerified bool                   `json:"email_verified" validate:"required"`
	PersonalInfo  UserPersonalInfoRecord `json:"personal_info" validate:"required"`
	Interests     []string               `json:"interests" validate:"required"`
}
