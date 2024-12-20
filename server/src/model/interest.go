package model

type Interest struct {
	Id       int    `json:"id"`
	Interest string `json:"name"`
}

type UserInterest struct {
	UserId       int    `json:"user_id" db:"user_id"`
	InterestId string 	`json:"interest_id" db:"interest_id"`
}
