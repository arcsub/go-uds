package uds

type List[T any] struct {
	Rows []T `json:"rows"`
}
