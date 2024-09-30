package model


type Pagination struct {
	NextOffset int `json:"next_offset,omitempty"`
	Limit      int `json:"limit" binding:"required"`
}

type PaginationResponse[T any] struct {
	Data       []T        `json:"data" binding:"required"`
	Pagination Pagination `json:"pagination" binding:"required"`
}

func CreatePaginationResponse[T any](data []T, limit, skip int, hasMore bool) PaginationResponse[T] {
	response := PaginationResponse[T]{
		Data: data,
		Pagination: Pagination{
			Limit: limit,
		},
	}

	if hasMore {
		response.Pagination.NextOffset = skip + limit
	}

	return response
}