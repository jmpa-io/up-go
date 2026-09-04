package up

// listOption is a struct used to help standardize the options used when listing
// any type of data from the API. It stores a name-value pair that represents a
// specific filtering or configuration option for list requests.
type listOption struct {
	name, value string // The name and value of the option.
}

// newListOption wraps the given name and value into a listOption struct.
// This helper function ensures the standardization of options used when
// listing any type of data from the API.
func newListOption(name, value string) listOption {
	return listOption{name: name, value: value}
}

// query returns the name and value of this option as query parameter key/value.
func (o listOption) query() (string, string) {
	return o.name, o.value
}

// queryProvider is implemented by all list-option wrapper types so setupQueries
// can iterate them without reflection.
type queryProvider interface {
	query() (string, string)
}
