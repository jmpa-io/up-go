package up

type Wrapper[T any] struct {
	Data  T     `json:"data"`
	Links Links `json:"links"`
}

type WrapperSlice[T any] struct {
	Data  []T   `json:"data"`
	Links Links `json:"links"`
}

type WrapperOmittable struct {
	Data  interface{} `json:"data,omitempty"`
	Links Links       `json:"links,omitempty"`
}
