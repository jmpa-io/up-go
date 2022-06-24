package up

type AccountObject struct {
	Data  DataObject         `json:"data"`
	Links RelatedLinksObject `json:"links"`
}
