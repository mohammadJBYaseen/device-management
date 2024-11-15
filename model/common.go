package model

type Page[T any] struct {
	PageNumber int32 `json:"page_number"`

	PageCount int32 `json:"page_count"`

	TotalCount int32 `json:"total_count"`

	TotalPages int32 `json:"total_pages"`

	Items []T `json:"items"`

	Sort Sort `json:"sort"`
}

type Sort struct {
	SortBy    string `json:"sort_by"`
	Direction string `json:"direction"`
}
