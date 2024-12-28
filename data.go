package up

type Data[T any, U any] struct {
	Object
	Attributes    T     `json:"attributes"`
	Relationships U     `json:"relationships"`
	Links         Links `json:"links"`
}
