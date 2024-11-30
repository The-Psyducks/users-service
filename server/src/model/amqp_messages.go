package model

type QueueMessage struct {
	MessageType string `json:"message_type"`
	Message     interface{} `json:"message"`
}

// LoginAttempt is a struct that represents a login attempt
type LoginAttempt struct {
	Succesfull	bool	`json:"was_succesfull"`
	UserId 		string	`json:"user_id"`
	Provider 	string	`json:"provider"`
	Timestamp 	string	`json:"timestamp"`
}

// NewRegistry is a struct that represents a new registry
type NewRegistry struct {
	RegistrationId 		string	`json:"registration_id"`
	Provider 			string	`json:"provider"`
	Timestamp 			string	`json:"timestamp"`
}

// UserBlocked is a struct that represents a user that has been blocked
type UserBlocked struct {
	UserId 		string	`json:"user_id"`
	Reason	 	string	`json:"reason"`
	Action	 	string	`json:"action"`
	Timestamp 	string	`json:"timestamp"`
}

// UserUnblocked is a struct that represents a user that has been unblocked
type UserUnblocked struct {
	UserId 		string	`json:"user_id"`
	Timestamp 	string	`json:"timestamp"`
}

// NewUser is a struct that represents a new user in the system
type NewUser struct {
	UserId 					string	`json:"user_id"`
	Location	 			string	`json:"location"`
	OldRegistrationId 		string	`json:"old_registration_id"`
	Timestamp 	string	`json:"timestamp"`
}