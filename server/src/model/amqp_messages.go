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