package page

// Page represents a paginated result set.
type Page[T any] struct {
	Records    []T   `json:"records"`
	Total      int64 `json:"total"`
	PageNum    int   `json:"pageNum"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}
