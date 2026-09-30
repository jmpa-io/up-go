package up

import "net/url"

// setupQueries converts a typed slice of list options into URL query
// parameters. Each option must implement queryProvider (via the embedded
// listOption). A default "page[size]" of 100 is added when not supplied.
func setupQueries[T queryProvider](options []T) url.Values {
	queries := make(url.Values)
	for _, o := range options {
		k, v := o.query()
		queries[k] = []string{v}
	}
	if _, ok := queries["page[size]"]; !ok {
		queries["page[size]"] = []string{"100"}
	}
	return queries
}
