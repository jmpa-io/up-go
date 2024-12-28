package up

type Links struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
	Next    string `json:"next,omitempty"`
	Prev    string `json:"prev,omitempty"`
}
