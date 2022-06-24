package up

type DataObject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// dataWrapper is a wrapper used when sending data to the API.
type dataWrapper struct {
	Data []interface{} `json:"data"`
}
