package model

type (
	Page[T any] struct {
		PageNumber int   `json:"page_number"`
		PageSize   int   `json:"page_size"`
		TotalCount int64 `json:"total_count"`
		TotalPages int   `json:"total_pages"`
		Items      []T   `json:"items"`
		Sort       Sort  `json:"sort"`
	}

	Sort struct {
		SortBy    string `json:"sort_by"`
		Direction string `json:"direction"`
	}

	SearchRequest struct {
		PageSize   int
		PageNumber int
		DeviceName string
		BrandName  string
		Sort       Sort
	}

	JsonPatch struct {
		Op    string `json:"op" binding:"required"`
		Value string `json:"value"`
		Path  string `json:"path" binding:"required"`
	}

	ApiError struct {
		Code string `json:"code,omitempty"`

		Message string `json:"message"`

		Domain string `json:"domain,omitempty"`

		DisplayMessage string `json:"display_message,omitempty"`
	}
)

func NewApiError(code string, message string, domain string, displayMessage string) *ApiError {
	return &ApiError{
		Code:           code,
		Message:        message,
		Domain:         domain,
		DisplayMessage: displayMessage,
	}
}
