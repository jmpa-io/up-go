package up

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

// Tags represents tags in Up
type Tags WrapperSlice[Object]

// wrapTags wraps the given tags in the data wrapper, ready to be sent to the API.
func wrapTags(tags []string) (wrappedTags Tags) {
	for _, t := range tags {
		wrappedTags.Data = append(wrappedTags.Data, Object{Type: "tags", ID: t})
	}
	return wrappedTags
}

// AddTagsToTransaction, using the given transaction id, adds the given tags to
// the transaction.
// https://developer.up.com.au/#post_transactions_transactionId_relationships_tags
func (c *Client) AddTagsToTransaction(ctx context.Context, id string, tags []string) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "AddTagsToTransaction")
	defer span.End()

	// add tags.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodPost,
		path:   fmt.Sprintf("/transactions/%s/relationships/tags", id),
		data:   wrapTags(tags),
	}, nil)
	return err
}

// RemoveTagsFromTransaction, using the given transaction id, removes the given
// tags from a transaction.
// https://developer.up.com.au/#delete_transactions_transactionId_relationships_tags
func (c *Client) RemoveTagsFromTransaction(ctx context.Context, id string, tags []string) error {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "RemoveTagsToTransaction")
	defer span.End()

	// remove tags.
	_, err := c.sender(newCtx, senderRequest{
		method: http.MethodDelete,
		path:   fmt.Sprintf("/transactions/%s/relationships/tags", id),
		data:   wrapTags(tags),
	}, nil)
	return err
}
