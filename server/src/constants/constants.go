package constants

// Validation constants
const (
	MinPasswordLength 	= 8
	MaxPasswordLength 	= 20
	MinUsernameLength 	= 4
	MaxUsernameLength 	= 20
	MinEmailLength    	= 3
	MaxEmailLength    	= 254
	MinFirstNameLength 	= 2
	MaxFirstNameLength 	= 100
	MinLastNameLength 	= 2
	MaxLastNameLength 	= 100
	MinInterests 		= 1
	MaxInterests 		= 100
)

// Resolver constants
const (
	LoginStep = "LOGIN"
	SignUpStep = "SIGN_UP"
	SessionStep  = "SESSION"
)

// onboarding constants
const (
	EmailVerificationStep = "EMAIL_VERIFICATION"
	PersonalInfoStep = "PERSONAL_INFO"
	InterestsStep = "INTERESTS"
	CompleteStep = "COMPLETE"
)

const (
	GoogleProvider = "GOOGLE"
	InternalProvider = "INTERNAL"
)

// Queue messages constants
const (
	LoginAttempt			= "LOGIN_ATTEMPT"
	UserBlocked 			= "USER_BLOCKED"
	UserUnblocked 			= "USER_UNBLOCKED"
	NewRegistry 			= "NEW_REGISTRY"
	NewUser 				= "NEW_USER"
)

const MaxPaginationLimit = 20