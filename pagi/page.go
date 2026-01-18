package pagi

type Params struct {
	Page  uint
	Limit uint
}

type Page[T any] struct {
	Data  T
	Page  uint
	Size  uint
	Total uint
}

type SortField struct {
	Ascend bool
	Field  string
}
