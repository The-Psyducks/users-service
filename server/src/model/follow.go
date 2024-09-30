package model

type UserPublicProfileWithFollowStatus struct {
	Follows bool              `json:"follows" binding:"required"`
	Profile UserPublicProfile `json:"profile" binding:"required"`
}

