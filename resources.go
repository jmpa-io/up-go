package up

type ResourceObject struct {
	ID    string         `json:"id"`
	Type  string         `json:"type"`
	Links SelfLinkObject `json:"links"`
}
