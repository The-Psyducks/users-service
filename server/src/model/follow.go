package model

type FollowUserPublicProfile struct {
	Follows bool              `json:"follows" binding:"required"`
	Profile UserPublicProfile `json:"profile" binding:"required"`
}

type Pagination struct {
	NextOffset	int	`json:"next_offset" binding:"required"`
	Limit		int	`json:"limit" binding:"required"`
}

type FollowersPaginationResponse struct {
	Followers []FollowUserPublicProfile `json:"data" binding:"required"`
	Pagination Pagination               `json:"pagination" binding:"required"`
}

type FollowingPaginationResponse struct {
	Following	[]FollowUserPublicProfile `json:"data" binding:"required"`
	Pagination	Pagination               `json:"pagination" binding:"required"`
}
