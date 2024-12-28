package up

type Wrapper[T any] struct {
	Data  T     `json:"data"`
	Links Links `json:"links"`
}

type WrapperSlice[T any] struct {
	Data  []T   `json:"data"`
	Links Links `json:"links"`
}

type WrapperOnlyLinks struct {
	Links Links `json:"links"`
}
