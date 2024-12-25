package up

type DataWrapper[T any] struct {
	Data T `json:"data"`
}

type SelfWrapper[T any] struct {
	Data  []T      `json:"data"`
	Links SelfLink `json:"links"`
}

type PaginationWrapper[T any] struct {
	Data  []T            `json:"data"`
	Links PaginationLink `json:"links"`
}
