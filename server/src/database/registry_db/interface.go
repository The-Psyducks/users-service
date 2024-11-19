package registry_db

import (
	"users-service/src/model"

	"github.com/google/uuid"
)

type RegistryDatabase interface {
	// CreateRegistryEntry creates a new registry entry with the given email
	CreateRegistryEntry(email string, identityProvider *string) (uuid.UUID, error)

	// GetRegistryEntry returns the registry entry with the given id
	GetRegistryEntry(id uuid.UUID) (model.RegistryEntry, error)

	// GetRegistryEntry returns the registry entry with the given id
	GetRegistryEntryByEmail(email string) (model.RegistryEntry, error)

	// AddPersonalInfoToRegistryEntry adds personal info to the registry entry with the given id
	AddPersonalInfoToRegistryEntry(id uuid.UUID, personalInfo model.UserPersonalInfoRecord) error

	// AddInterestsToRegistryEntry adds interests to the registry entry with the given id
	AddInterestsToRegistryEntry(id uuid.UUID, interests []string) error

	// SetEmailVerificationPin sets the email verification pin of the registry entry with the given id
	SetEmailVerificationPin(id uuid.UUID, code string) error

	// GetEmailVerificationPin returns the email verification pin of the registry entry with the given id
	GetEmailVerificationPin(id uuid.UUID) (string, error)

	// VerificateEmail notifies that the email of the registry entry with the given id has been verified
	VerifyEmail(id uuid.UUID) error

	// CheckIfRegistryEntryExists checks if a registry entry with the given id exists
	CheckIfRegistryEntryExists(id uuid.UUID) (bool, error)

	// CheckIfRegistryEntryExistsByEmail checks if a registry entry with the given email exists
	CheckIfRegistryEntryExistsByEmail(email string) (bool, error)

	// DeleteRegistryEntry deletes the registry entry with the given id
	DeleteRegistryEntry(id uuid.UUID) error

	// GetRegistrySummaryMetrics returns the summary metrics of the registry
	GetRegistrySummaryMetrics() (*model.RegistrationSummaryMetrics, error)
}
