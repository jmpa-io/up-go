package up

import (
	"net/url"
	"reflect"
)

// setupQueries takes a slice of options with "name" and "value" fields and
// converts them into URL query parameters. It uses reflection to handle
// different slice types and ensures "page[size]" defaults to "100" if not set.
func setupQueries(options interface{}) url.Values {
	queries := make(url.Values)

	// add specific fields if the given options are a slice.
	// TODO: any validation on the types of options here?
	v := reflect.ValueOf(options)
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			// add the "name" field as the key & the "value" field as the value.
			queries[item.FieldByName("name").String()] = []string{
				item.FieldByName("value").String(),
			}
		}
	}

	// add default "page[size]" if it's not already set.
	if _, ok := queries["page[size]"]; !ok {
		queries["page[size]"] = []string{"100"}
	}

	return queries
}
