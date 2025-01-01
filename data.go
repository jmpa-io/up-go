package up

// Data is a generic struct that represents a resource returned from the API.
// It contains the resource's attributes, relationships, and links. The struct
// is parameterized with two types: T for the attributes and U for the
// relationships, allowing it to be used for various types of resources.
type Data[T any, U any] struct {
	Object
	Attributes    T     `json:"attributes"`
	Relationships U     `json:"relationships"`
	Links         Links `json:"links"`
}
