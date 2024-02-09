package persistence

type Pagination[T any] struct {
	Page   int64 `json:"page"`
	Size   int64 `json:"size"`
	Result T     `json:"result"`
}

func NewPagination[T any](page int64, size int64, res T) *Pagination[T] {
	return &Pagination[T]{
		Page:   page,
		Size:   size,
		Result: res,
	}
}
