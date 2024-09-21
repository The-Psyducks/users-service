package model


type FollowUserPublicProfile struct {
	Follows    bool					`json:"follows" binding:"required"`
	Profile    UserPublicProfile	`json:"profile" binding:"required"`
}

type FollowersResponse struct {
	Followers []FollowUserPublicProfile	`json:"followers" binding:"required"`
}

type FollowingResponse struct {
	Following []FollowUserPublicProfile	`json:"following" binding:"required"`
}
