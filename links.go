package up

// Links represents a collection of URLs for a resource returned from the API.
// It can contain links to the resource itself, related resources, and pagination
// links (next and previous pages). Only relevant fields will be populated based
// on the context, and fields can be omitted when not needed.
type Links struct {

	// Self represents the URL of the current resource.
	Self string `json:"self,omitempty"`

	// Related represents a URL to a related resource.
	Related string `json:"related,omitempty"`

	// Next represents the URL for the next page of results, if paginating.
	Next string `json:"next,omitempty"`

	// Prev represents the URL for the previous page of results, if paginating.
	Prev string `json:"prev,omitempty"`
}
