package up

import (
	"fmt"
	"net/http"
)

type TagObject struct {
	Data  DataObject     `json:"data"`
	Links SelfLinkObject `json:"links"`
}

// wrapTags wraps the given tags in the data wrapper, ready to be sent to the API.
func wrapTags(tags []string) (d dataWrapper) {
	for _, t := range tags {
		w := struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		}{
			Type: "tags",
			ID:   t,
		}
		d.Data = append(d.Data, w)
	}
	return d
}

// AddTagsToTransaction, using the given transaction id, adds the given tags to
// the transaction.
// https://developer.up.com.au/#post_transactions_transactionId_relationships_tags
func (c *Client) AddTagsToTransaction(id string, tags []string) error {
	cr := clientRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transactions/%s/relationships/tags", id),
		data:   wrapTags(tags),
	}
	_, err := c.sender(cr, nil)
	return err
}

// RemoveTagsFromTransaction, using the given transaction id, removes the given
// tags from a transaction.
// https://developer.up.com.au/#delete_transactions_transactionId_relationships_tags
func (c *Client) RemoveTagsFromTransaction(id string, tags []string) error {
	cr := clientRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/transactions/%s/relationships/tags", id),
		data:   wrapTags(tags),
	}
	_, err := c.sender(cr, nil)
	return err
}
