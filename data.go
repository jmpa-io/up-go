package up

type Data[T any, U any] struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Attributes    T        `json:"attributes"`
	Relationships U        `json:"relationships"`
	Links         SelfLink `json:"links"`
}
