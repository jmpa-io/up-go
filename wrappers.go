package up

// Wrapper is a generic struct that wraps data of type T under the `Data` field,
// along with additional metadata under the `Links` field. It is commonly used
// to package data sent to or received from the API.
type Wrapper[T any] struct {
	Data  T     `json:"data"`
	Links Links `json:"links"`
}

// WrapperSlice is a generic struct that wraps a slice of type T under the
// `Data` field, along with additional metadata under the `Links` field. It is
// commonly used to package data sent to or received from the API.
type WrapperSlice[T any] struct {
	Data  []T   `json:"data"`  // A slice of data to be wrapped.
	Links Links `json:"links"` // Metadata or related links.
}

// WrapperOmittable is a struct that wraps data under the `Data` field, along
// with additional metadata under the `Links` field. Both fields are optional,
// and will be omitted from JSON output if empty.
type WrapperOmittable struct {
	Data  interface{} `json:"data,omitempty"`
	Links Links       `json:"links,omitempty"`
}
